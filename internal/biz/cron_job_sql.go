package biz

import (
	"context"
	"cron/internal/basic/config"
	"cron/internal/basic/db"
	"cron/internal/models"
	"cron/internal/pb"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"gorm.io/gorm"
	"strings"
	"time"
)

// mysql 命令执行
func (job *CronJob) sqlMysql(ctx context.Context, r *pb.CronSql) (resp []byte, err error) {
	conf := config.DataBaseConf{
		Source: fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=false&loc=Local",
			r.Source.Username, r.Source.Password, r.Source.Hostname, r.Source.Port, r.Source.Database),
	}
	_db := db.Conn(conf).WithContext(ctx)
	if _db.Error != nil {
		return nil, _db.Error
	}

	// 执行sql
	res := make([]string, len(r.Statement))
	for i, sql := range r.Statement {
		str, err := job.sqlMysqlItem(_db, sql)
		res[i] = str
		if err != nil && r.ErrAction != models.SqlErrActionProceed {
			break
		}
	}
	return jsoniter.Marshal(strings.Join(res, "\r\n"))
}

func (job *CronJob) sqlMysqlItem(_db *gorm.DB, sql string) (res string, err error) {
	startTime := time.Now()
	defer func() {
		endTime := time.Now()
		status := "成功"
		info := "ok."
		if err != nil {
			status = "失败"
			info = err.Error()
		}
		res = fmt.Sprintf("[%s]	%v	|	%s	|	%s", startTime.Format(time.RFC3339), endTime.Unix()-startTime.Unix(), status, info)
	}()

	err = _db.Exec(sql).Error
	return res, err
}
