package biz

import (
	"context"
	"cron/internal/basic/tracing"
	"cron/internal/models"
	"log"
	"os"
	"testing"
)

func init() {
	err := os.Chdir("E:/WorkApps/go/src/incron/")
	if err != nil {
		log.Fatal(err)
	}
	// 日志写入
	go tracing.MysqlCollectorListen()
}

func TestJenKins(t *testing.T) {
	conf := &models.CronConfig{
		Id:           0,
		Env:          "",
		EntryId:      0,
		Type:         0,
		Name:         "",
		Spec:         "",
		Protocol:     0,
		Command:      []byte(`{"cmd": "", "rpc": {"addr": "", "body": "", "proto": "", "action": "", "header": [], "method": "GRPC", "actions": []}, "sql": {"driver": "mysql", "source": {"id": 0, "port": "", "title": "", "database": "", "hostname": "", "password": "", "username": ""}, "interval": 0, "statement": [], "err_action": 1, "err_action_name": ""}, "http": {"url": "", "body": "", "header": [{"key": "", "value": ""}], "method": "GET"}, "jenkins": {"name": "card", "params": [{"key": "", "value": ""}], "source": {"id": 10}}}`),
		Remark:       "",
		Status:       0,
		StatusRemark: "",
		StatusDt:     "",
		UpdateDt:     "",
		CreateDt:     "",
		MsgSet:       nil,
	}

	ctx := context.Background()
	job := NewJobConfig(conf)

	job.jenkins(ctx, job.commandParse.Jenkins)
}

func TestJenkins2(t *testing.T) {

}
