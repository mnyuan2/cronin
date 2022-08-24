package biz

import (
	"log"
	"testing"
	"time"
)

func TestInitTask(t *testing.T) {
	task := NewTaskService()
	task.Init()

	time.Sleep(time.Second * 20)
	log.Println("任务执行完毕...")
}
