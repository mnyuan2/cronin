package biz

import (
	"context"
	"cron/internal/basic/config"
	"cron/internal/basic/db"
	"cron/internal/basic/enum"
	"cron/internal/basic/errs"
	"cron/internal/basic/tracing"
	"cron/internal/data"
	"cron/internal/models"
	"cron/internal/pb"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
	"time"
)

// mysql 命令执行
func (job *JobConfig) sqlMysql(ctx context.Context, r *pb.CronSql) (err errs.Errs) {
	ctx, span := job.tracer.Start(ctx, "exec-mysql")
	defer func() {
		if err != nil {
			span.SetStatus(tracing.StatusError, err.Desc())
			span.AddEvent("执行错误", trace.WithAttributes(attribute.String("error.object", err.Error())))
		} else {
			span.SetStatus(tracing.StatusOk, "")
		}
		span.End()
	}()
	b, _ := jsoniter.Marshal(r)
	span.AddEvent("sql_set", trace.WithAttributes(
		attribute.String("sql_set", string(b)),
	))

	source, er := data.NewCronSettingData(ctx).GetSourceOne(job.conf.Env, r.Source.Id)
	if er != nil {
		return errs.New(er, "连接配置异常")
	}
	s := &pb.CronSqlSource{}
	if er = jsoniter.UnmarshalFromString(source.Content, s); er != nil {
		return errs.New(er, "连接配置解析异常")
	}

	password, er := models.SqlSourceDecode(s.Password)
	if er != nil {
		return errs.New(er, "密码异常")
	}
	conf := &config.MysqlSource{
		Hostname: s.Hostname,
		Database: s.Database,
		Username: s.Username,
		Password: password,
		Port:     s.Port,
	}
	_db := db.Conn(conf).WithContext(ctx)

	if _db.Error != nil {
		span.AddEvent("连接失败", trace.WithAttributes(
			attribute.String("dsn", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=false&loc=Local",
				conf.Username, conf.Password, conf.Hostname, conf.Port, conf.Database))),
		)
		return errs.New(_db.Error, "连接失败")
	}

	// 执行sql
	// 此处为局部错误，大框架是完成的，错误不必返回
	err = job.sqlMysqlExec(r, _db, r.Statement)
	if err != nil {
		// 执行告警推送
		go job.messagePush(ctx, enum.StatusDisable, err.Desc(), []byte(err.Error()), 0)
	}
	return err
}

func (job *JobConfig) sqlMysqlExec(r *pb.CronSql, _db *gorm.DB, statement []string) errs.Errs {
	var tx *gorm.DB
	if r.ErrAction == models.SqlErrActionRollback {
		tx = _db.Begin()
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

func (job *JobConfig) sqlMysqlItem(r *pb.CronSql, _db *gorm.DB, sql string) (err error) {
	ctx, span := job.tracer.Start(_db.Statement.Context, "sql-item")
	span.AddEvent("", trace.WithAttributes(attribute.String("sql", sql)))
	defer span.End()

	resp := _db.Exec(sql)
	if resp.Error != nil {
		err = resp.Error
		span.SetStatus(codes.Error, "执行失败")
		span.AddEvent("error", trace.WithAttributes(attribute.String("error.object", err.Error())))
		if r.ErrAction == models.SqlErrActionProceed {
			go job.messagePush(ctx, enum.StatusDisable, "错误跳过继续", []byte(err.Error()), 0)
		}
	} else {
		span.SetStatus(codes.Ok, "成功")
		span.AddEvent("", trace.WithAttributes(attribute.Int64("rows affected", resp.RowsAffected)))
	}
	return err
}
