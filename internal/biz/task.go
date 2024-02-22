package biz

import (
	"context"
	"cron/internal/basic/config"
	"cron/internal/basic/conv"
	"cron/internal/basic/db"
	"cron/internal/basic/enum"
	"cron/internal/basic/errs"
	"cron/internal/data"
	"cron/internal/models"
	"cron/internal/pb"
	"encoding/base64"
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/robfig/cron/v3"
	"net/http"
	"time"
)

type TaskService struct {
	cron *cron.Cron // 任务计划 组件
	conf config.Main
}

func NewTaskService(conf config.Main) *TaskService {
	return &TaskService{
		cron: cronRun,
		conf: conf,
	}
}

// Init 初始化任务
func (dm *TaskService) Init() (err error) {
	pageSize, total := 500, int64(500)
	cronDb := data.NewCronConfigData(context.Background())
	for page := 1; total >= int64(pageSize*page); page++ {
		list := []*models.CronConfig{}
		w := db.NewWhere().Eq("status", enum.StatusActive)
		total, err = cronDb.ListPage(w, page, pageSize, &list)
		if err != nil {
			panic(fmt.Sprintf("任务配置读取异常：%s", err.Error()))
		}
		for _, conf := range list {
			// 启用成功，更新任务id；启动失败，置空任务id
			if err := dm.Add(conf); err != nil {
				conf.EntryId = 0
				conf.Status = models.ConfigStatusError
				cronDb.ChangeStatus(conf, "初始化注册错误")
			}
		}
	}

	// 系统内置任务
	dm.Add(dm.sysLogRetentionConf())

	return nil
}

// 注册任务监视器
func (dm *TaskService) sysRegisterMonitor() {
	// 检测系统任务是否正常
	// 系统任务包含：定时删除日志、注册任务有效性检查。
	//cronDb := data.NewCronConfigData(context.Background())
	//w := db.NewWhere().Eq("status", enum.StatusActive)
	//total, err = cronDb.GetList(w, 1, pageSize, &list)
}

// 添加任务
func (dm *TaskService) Add(conf *models.CronConfig) error {
	if conf == nil {
		return errors.New("未指定任务")
	}
	j := NewCronJob(conf)
	if conf.Type == models.TypeOnce {
		return dm.addOnce(j)
	}
	id, err := dm.cron.AddJob(conf.Spec, j)
	if err != nil {
		g := models.NewErrorCronLog(conf, "", errs.New(err, "任务启动失败"), time.Now())
		data.NewCronLogData(context.Background()).Add(g)
		return err
	}
	conf.EntryId = int(id)
	return nil
}

// 添加单次任务
func (dm *TaskService) addOnce(j *JobConfig) error {
	s, err := NewScheduleOnce(j.conf.Spec)
	if err != nil {
		return err
	}
	id := dm.cron.Schedule(s, j)
	j.conf.EntryId = int(id)
	return nil
}

// 删除任务
func (dm *TaskService) Del(conf *models.CronConfig) {
	if conf.EntryId == 0 {
		return
	}
	dm.cron.Remove(cron.EntryID(conf.EntryId))
}

// sysLogDurationConf 内置任务，日志删除
func (dm *TaskService) sysLogRetentionConf() *models.CronConfig {
	retention := dm.conf.Task.LogRetention
	if retention == "" {
		return nil
	}
	re, err := time.ParseDuration(retention)
	if err != nil {
		panic(fmt.Sprintf("log_retention 日志存续配置有误, %s", err.Error()))
	} else if re.Hours() < 24 {
		panic("log_retention 日志存续不得小于24h")
	}

	var sysLogRetention = &models.CronConfig{
		Id:       -1,
		Name:     "日志留存时间",
		Spec:     "0 0 5 * * *", // 每天5点执行
		Protocol: models.ProtocolHttp,
		Status:   enum.StatusActive,
		Remark:   "系统内置任务",
		CreateDt: time.Now().Format(conv.FORMAT_DATETIME),
		UpdateDt: time.Now().Format(conv.FORMAT_DATETIME),
	}
	cmd := &pb.CronConfigCommand{
		Http: &pb.CronHttp{
			Method: http.MethodPost,
			Url:    dm.conf.Http.Local() + "/log/del",
			Body:   fmt.Sprintf(`{"retention":"%s"}`, dm.conf.Task.LogRetention),
		},
	}
	// 启用了账号时，构建token header
	if config.MainConf().User.AdminAccount != "" {
		s := base64.StdEncoding.EncodeToString([]byte(config.MainConf().User.AdminAccount + ":" + config.MainConf().User.AdminPassword))
		cmd.Http.Header = []*pb.KvItem{
			{Key: "Authorization", Value: "Basic " + s},
		}
	}
	sysLogRetention.Command, _ = jsoniter.Marshal(cmd)
	return sysLogRetention
}
