package biz

import (
	"context"
	"cron/internal/basic/config"
	"cron/internal/basic/db"
	"cron/internal/basic/enum"
	"cron/internal/basic/errs"
	"cron/internal/data"
	"cron/internal/models"
	"cron/internal/pb"
	jsoniter "github.com/json-iterator/go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

// mysql 命令执行
func (job *CronJob) sqlMysql(ctx context.Context, r *pb.CronSql) (err error) {
	ctx, span := job.tracer.Start(ctx, "sqlMysql")
	span.SetStatus(codes.Ok, "")
	defer func() {
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	source, err := data.NewCronSettingData(ctx).GetSqlSourceOne(job.conf.Env, r.Source.Id)
	if err != nil {
		return errs.New(err, "连接配置异常")
	}
	s := &pb.CronSqlSource{}
	if err = jsoniter.UnmarshalFromString(source.Content, s); err != nil {
		return errs.New(err, "连接配置解析异常")
	}

	password, err := models.SqlSourceDecode(s.Password)
	if err != nil {
		return errs.New(err, "密码异常")
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
		return errs.New(_db.Error, "连接失败")
	}

	// 执行sql
	// 此处为局部错误，大框架是完成的，错误不必返回
	err = job.sqlMysqlExec(_db, r.ErrAction, r.Statement)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		// 执行告警推送
		job.messagePush(ctx, enum.StatusDisable, err.Error(), nil, 0)
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
	_, span := job.tracer.Start(_db.Statement.Context, "sqlMysqlItem")
	span.AddEvent("x", trace.WithAttributes(attribute.String("sql", "sql")))
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
