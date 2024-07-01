package biz

import (
	"context"
	"cron/internal/basic/config"
	"cron/internal/basic/db"
	"cron/internal/basic/tracing"
	"cron/internal/models"
	"fmt"
	"log"
	"os"
	"testing"
	"time"
)

func init() {
	os.Chdir("../../") // 设置运行根目录
}

func TestInitTask(t *testing.T) {
	task := NewTaskService(config.MainConf())
	task.Init()

	time.Sleep(time.Second * 20)
	log.Println("任务执行完毕...")
}

func TestTaskDemo1(t *testing.T) {
	// 日志写入
	go tracing.MysqlCollectorListen()
	task := NewTaskService(config.MainConf())
	_db := db.New(context.Background())

	conf := &models.CronConfig{}
	_db.Where("id=?", 133).Find(conf)
	if conf.Id == 0 {
		t.Fatal("未找到任务")
	}
	conf.Spec = "*/3 * * * * *"

	task.AddConfig(conf)

	time.Sleep(120 * time.Second)
	t.Log("end...")
}

func TestParseTime(t *testing.T) {
	ti := time.Now()

	fmt.Println(ti.Unix(), ti.UnixMilli(), ti.UnixMicro())
	spec := "2024-11-02 00:25:24"

	spec = "59 59 23 * * */7" // 每周 7 23:59:59
	//spec = "0 0 0 * * 1"      // 每周 1 0:0:0    这里有点奇怪为什么周7和周1的写法不一样！
	su, err := secondParser.Parse(spec)
	if err != nil {
		t.Fatal(err)
	}
	nti := su.Next(ti)
	fmt.Println(su, "\n	", nti, "\n	", err)
	for i := 0; i < 3; i++ {
		nnti := su.Next(nti)
		fmt.Println("	", nnti)
		nti = nnti
	}

	d, err := time.ParseDuration("24h")
	t1 := time.Now().Add(-d)
	fmt.Println(err, d, d.Hours(), d.Hours() < 24, t1)
}
