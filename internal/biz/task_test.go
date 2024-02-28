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
	ti := time.Now()
	spec := "2022-11-02 00:25:24"
	spec = "59 59 23 * * */7" // 每周 7 23:59:59
	su, err := secondParser.Parse(spec)
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
