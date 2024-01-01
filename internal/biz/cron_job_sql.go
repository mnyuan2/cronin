package biz

import (
	"context"
	"cron/internal/basic/config"
	"cron/internal/basic/db"
	"cron/internal/data"
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
	source, err := data.NewCronSettingData(ctx).GetSqlSourceOne(job.conf.Env, r.Source.Id)
	if err != nil {
		return nil, fmt.Errorf("连接配置异常 %w", err)
	}
	s := &pb.CronSqlSource{}
	if err = jsoniter.UnmarshalFromString(source.Content, s); err != nil {
		return nil, fmt.Errorf("连接配置解析异常 %w", err)
	}

	password, err := models.SqlSourceDecode(s.Password)
	if err != nil {
		return nil, fmt.Errorf("密码异常,%w", err)
	}
	conf := config.DataBaseConf{
		Source: fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=false&loc=Local",
			s.Username, password, s.Hostname, s.Port, s.Database),
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
	return []byte(strings.Join(res, `
`)), nil
}

func (job *CronJob) sqlMysqlItem(_db *gorm.DB, sql string) (res string, err error) {
	startTime := time.Now()
	info := ""
	defer func() {
		status := "成功"
		if err != nil {
			status = "失败"
			info = err.Error()
		}
		res = fmt.Sprintf("[%s]	%vs	|	%s	|	%s", startTime.Format(time.RFC3339), time.Since(startTime).Seconds(), status, info)
	}()

	resp := _db.Exec(sql)
	if resp.Error != nil {
		err = resp.Error
	} else {
		info = fmt.Sprintf("rows affected: %v", resp.RowsAffected)
	}
	return res, err
}
