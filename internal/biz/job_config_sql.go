package biz

import (
	"bytes"
	"context"
	"cron/internal/basic/config"
	"cron/internal/basic/db"
	"cron/internal/basic/enum"
	"cron/internal/basic/errs"
	"cron/internal/basic/tracing"
	"cron/internal/basic/util"
	"cron/internal/biz/dtos"
	"cron/internal/data"
	"cron/internal/models"
	"cron/internal/pb"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

// sql 执行函数
func (job *JobConfig) sqlFunc(ctx context.Context) (err errs.Errs) {
	switch job.commandParse.Sql.Driver {
	case enum.SqlDriverMysql, enum.SqlDriverClickhouse:
		return job.sql(ctx, job.commandParse.Sql)
	default:
		return errs.New(nil, fmt.Sprintf("未支持的sql 驱动 %s", job.commandParse.Sql.Driver), errs.SysError)
	}
}

// sql 语句执行
func (job *JobConfig) sql(ctx context.Context, r *pb.CronSql) (err errs.Errs) {
	ctx, span := job.tracer.Start(ctx, "exec-"+r.Driver)
	defer func() {
		if er := util.PanicInfo(recover()); er != "" {
			span.SetStatus(tracing.StatusError, "执行异常")
			span.AddEvent("error", trace.WithAttributes(attribute.String("error.panic", er)))
		} else if err != nil {
			span.SetStatus(tracing.StatusError, err.Desc())
			span.AddEvent("执行错误", trace.WithAttributes(attribute.String("error.object", err.Error())))
		} else {
			span.SetStatus(tracing.StatusOk, "")
		}
		span.End()
	}()
	span.SetAttributes(attribute.String("driver", r.Driver))
	//b, _ := jsoniter.Marshal(r)
	//span.AddEvent("sql_set", trace.WithAttributes(
	//attribute.String("sql_set", string(b)),
	//))

	source, er := data.NewCronSettingData(ctx).GetSourceOne(job.conf.Env, r.Source.Id)
	if er != nil {
		return errs.New(er, "连接配置异常")
	}
	s, er := dtos.ParseSource(source)
	if er != nil {
		return errs.New(er, "连接配置解析异常")
	}

	password, er := models.SqlSourceDecode(s.Sql.Password)
	if er != nil {
		return errs.New(er, "密码异常")
	}
	conf := &config.MysqlSource{
		Hostname: s.Sql.Hostname,
		Database: s.Sql.Database,
		Username: s.Sql.Username,
		Password: password,
		Port:     s.Sql.Port,
	}
	_db := &gorm.DB{}
	if r.Driver == enum.SqlDriverMysql {
		_db = db.Conn(conf).WithContext(ctx)
	} else {
		_db = db.ConnClickhouse(conf).WithContext(ctx)
	}

	if _db.Error != nil {
		span.AddEvent("连接失败", trace.WithAttributes(
			attribute.String("dsn", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=false&loc=Local",
				conf.Username, conf.Password, conf.Hostname, conf.Port, conf.Database))),
		)
		return errs.New(_db.Error, "连接失败")
	}
	_db.Logger = logger.Discard // 取消日志打印

	statement := []*pb.KvItem{} // value.具体sql、key.描述备注

	// 仅支持与origin一致的sql语句
	for _, item := range r.Statement {
		if item.Type != r.Origin {
			continue
		}
		switch item.Type {
		case enum.SqlStatementSourceGit:
			files, err := job.getGitFile(ctx, item.Git)
			if err != nil {
				return err
			}
			for _, file := range files {
				list := [][]byte{}
				if item.IsBatch == enum.BoolYes {
					list = bytes.Split(file.Byte, []byte(";"))
				} else if item.IsBatch == enum.BoolNot {
					list = [][]byte{0: file.Byte}
				} else {
					return errs.New(nil, "git 批量解析标识未填写，请重新编辑后提交")
				}
				for _, item := range list {
					s := bytes.TrimSpace(item)
					if s != nil {
						statement = append(statement, &pb.KvItem{
							Key:   file.Name,
							Value: string(s),
						})
					}
				}
			}
		case enum.SqlStatementSourceLocal:
			key := util.ParseSqlTypeName(item.Local)
			statement = append(statement, &pb.KvItem{Key: key, Value: item.Local})
		default:
			return errs.New(_db.Error, "sql 语句来源异常")
		}
	}

	// 执行sql
	// 此处为局部错误，大框架是完成的，错误不必返回
	err = job.sqlMysqlExec(r, _db, statement)
	if err != nil {
		// 执行告警推送
		go job.messagePush(ctx, &dtos.MsgPushRequest{
			Status:     enum.StatusDisable,
			StatusDesc: err.Desc(),
			Body:       []byte(err.Error()),
			Duration:   0,
			RetryNum:   0,
		})
	}
	return err
}

func (job *JobConfig) sqlMysqlExec(r *pb.CronSql, _db *gorm.DB, statement []*pb.KvItem) errs.Errs {
	var tx *gorm.DB
	if r.ErrAction == models.SqlErrActionRollback {
		tx = _db.Begin()
		r.Interval = 0
	} else {
		tx = _db
	}

	for i, sql := range statement {
		if i > 0 && r.Interval > 0 {
			time.Sleep(time.Second * time.Duration(r.Interval)) // 间隔秒
		}
		err := job.sqlMysqlItem(r, tx, sql)
		if err != nil {
			if r.ErrAction == models.SqlErrActionAbort { // 终止
				return errs.New(err, "错误终止")
			} else if r.ErrAction == models.SqlErrActionRollback { // 回滚
				tx.Rollback()
				return errs.New(err, "错误回滚")
			} else if r.ErrAction == models.SqlErrActionProceed { // 继续

			}
		}
	}

	if r.ErrAction == models.SqlErrActionRollback {
		tx = tx.Commit()
	}
	return nil
}

func (job *JobConfig) sqlMysqlItem(r *pb.CronSql, _db *gorm.DB, item *pb.KvItem) (err error) {
	ctx, span := job.tracer.Start(_db.Statement.Context, "sql-item")
	span.AddEvent("", trace.WithAttributes(attribute.String("sql", item.Value)))
	defer span.End()

	resp := &gorm.DB{}
	if item.Key == "SELECT" {
		dataStr := ""
		data := []map[string]any{}
		resp = _db.Raw(item.Value).Scan(&data)
		if len(data) > 0 {
			dataStr, _ = jsoniter.MarshalToString(data)
		}
		span.SetAttributes(attribute.String("type", item.Key))
		span.AddEvent("", trace.WithAttributes(attribute.String("rows", dataStr)))
	} else {
		resp = _db.Exec(item.Value)
	}

	if resp.Error != nil {
		err = resp.Error
		span.SetStatus(codes.Error, "执行失败")
		span.AddEvent("error", trace.WithAttributes(
			attribute.String("error.object", err.Error()),
			attribute.String("remark", models.SqlErrActionMap[r.ErrAction]),
		))
		if r.ErrAction == models.SqlErrActionProceed {
			go job.messagePush(ctx, &dtos.MsgPushRequest{
				Status:     enum.StatusDisable,
				StatusDesc: "错误跳过继续",
				Body:       []byte(err.Error()),
				Duration:   0,
				RetryNum:   0,
			})
		}
	} else {
		span.SetStatus(codes.Ok, "成功")
		span.AddEvent("", trace.WithAttributes(
			attribute.Int64("rows affected", resp.RowsAffected),
		))
	}
	return err
}
