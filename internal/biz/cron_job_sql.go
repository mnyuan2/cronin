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
)

// mysql 命令执行
func (job *CronJob) sqlMysql(ctx context.Context, r *pb.CronSql) (err errs.Errs) {
	ctx, span := job.tracer.Start(ctx, "exec-mysql")
	defer func() {
		if err != nil {
			span.SetStatus(tracing.Error, err.Desc())
			span.AddEvent("x", trace.WithAttributes(attribute.String("error.object", err.Error())))
		} else {
			span.SetStatus(tracing.Ok, "")
		}
		span.End()
	}()

	source, er := data.NewCronSettingData(ctx).GetSqlSourceOne(job.conf.Env, r.Source.Id)
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
		span.AddEvent("x", trace.WithAttributes(
			attribute.String("dsn", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=false&loc=Local",
				conf.Username, conf.Password, conf.Hostname, conf.Port, conf.Database))),
		)
		return errs.New(_db.Error, "连接失败")
	}

	// 执行sql
	// 此处为局部错误，大框架是完成的，错误不必返回
	er = job.sqlMysqlExec(_db, r.ErrAction, r.Statement)
	if er != nil {
		span.SetStatus(codes.Error, er.Error())
		// 执行告警推送
		job.messagePush(ctx, enum.StatusDisable, er.Error(), nil, 0)
	}
	return nil
}

func (job *CronJob) sqlMysqlExec(_db *gorm.DB, errAction int, statement []string) (err error) {
	var tx *gorm.DB
	if errAction == models.SqlErrActionRollback {
		tx = _db.Begin()
	} else {
		tx = _db
	}

	for _, sql := range statement {
		err = job.sqlMysqlItem(tx, sql)
		if err != nil {
			if errAction == models.SqlErrActionAbort { // 终止
				return errs.New(err, "错误终止")
			} else if errAction == models.SqlErrActionRollback { // 回滚
				tx.Rollback()
				return errs.New(err, "错误回滚")
			} else if errAction == models.SqlErrActionProceed { // 继续

			}
		}
	}

	if errAction == models.SqlErrActionRollback {
		tx = tx.Commit()
	}
	return nil
}

func (job *CronJob) sqlMysqlItem(_db *gorm.DB, sql string) (err error) {
	_, span := job.tracer.Start(_db.Statement.Context, "sql-item")
	span.AddEvent("x", trace.WithAttributes(attribute.String("sql", sql)))
	defer span.End()

	resp := _db.Exec(sql)
	if resp.Error != nil {
		err = resp.Error
		span.SetStatus(codes.Error, err.Error())
	} else {
		span.SetStatus(codes.Ok, "成功")
		span.AddEvent("x", trace.WithAttributes(attribute.Int64("rows affected", resp.RowsAffected)))
	}
	return err
}
