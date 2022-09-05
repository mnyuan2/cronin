package biz

import (
	"cron/internal/basic/config"
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

func TestParseTime(t *testing.T) {
	su, err := secondParser.Parse("* 0/5 * * * ?")
	fmt.Println(su, err)

	d, err := time.ParseDuration("24h")
	t1 := time.Now().Add(-d)
	fmt.Println(err, d, d.Hours(), d.Hours() < 24, t1)
}
