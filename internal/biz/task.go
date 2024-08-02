package biz

import (
	"context"
	"cron/internal/basic/auth"
	"cron/internal/basic/config"
	"cron/internal/basic/conv"
	"cron/internal/basic/db"
	"cron/internal/basic/enum"
	"cron/internal/basic/sse"
	"cron/internal/basic/tracing"
	"cron/internal/biz/dtos"
	"cron/internal/data"
	"cron/internal/models"
	"cron/internal/pb"
	"encoding/base64"
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/robfig/cron/v3"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"log"
	"net/http"
	"sync"
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
	logDb := data.NewCronChangeLogData(context.Background())
	for page := 1; total >= int64(pageSize*page); page++ {
		list := []*models.CronConfig{}
		w := db.NewWhere().Eq("status", enum.StatusActive).In("type", []int{models.TypeCycle, models.TypeOnce})
		total, err = cronDb.ListPage(w, page, pageSize, &list)
		if err != nil {
			panic(fmt.Sprintf("任务配置读取异常：%s", err.Error()))
		}
		for _, conf := range list {
			g := data.NewChangeLogHandle(&auth.UserToken{UserName: "系统"}).SetType(models.LogTypeUpdateSys).OldConfig(*conf)
			// 启用成功，更新任务id；启动失败，置空任务id
			if err := dm.AddConfig(conf); err != nil {
				conf.EntryId = 0
				conf.Status = models.ConfigStatusError
				cronDb.ChangeStatus(conf, "初始化 "+err.Error())
			} else {
				cronDb.SetEntryId(conf)
			}
			logDb.Write(g.NewConfig(*conf))
		}
	}

	pageSize, total = 500, int64(500)
	pipeineData := data.NewCronPipelineData(context.Background())
	for page := 1; total >= int64(pageSize*page); page++ {
		list := []*models.CronPipeline{}
		w := db.NewWhere().Eq("status", enum.StatusActive)
		total, err = pipeineData.ListPage(w, page, pageSize, &list)
		if err != nil {
			panic(fmt.Sprintf("任务配置读取异常：%s", err.Error()))
		}
		for _, conf := range list {
			g := data.NewChangeLogHandle(&auth.UserToken{UserName: "系统"}).SetType(models.LogTypeUpdateSys).OldPipeline(*conf)
			if err := dm.AddPipeline(conf); err != nil {
				conf.EntryId = 0
				conf.Status = models.ConfigStatusError
				pipeineData.ChangeStatus(conf, "初始化 "+err.Error())
			} else {
				pipeineData.SetEntryId(conf)
			}
			logDb.Write(g.NewPipeline(*conf))
		}
	}

	// 系统内置任务
	dm.AddConfig(dm.sysLogRetentionConf())

	return nil
}

// 注册任务监视器
func (dm *TaskService) RegisterMonitor() {
	// 检测系统任务是否正常
	// 系统任务包含：定时删除日志、注册任务有效性检查。
	//cronDb := data.NewCronConfigData(context.Background())
	//w := db.NewWhere().Eq("status", enum.StatusActive)
	//total, err = cronDb.GetList(w, 1, pageSize, &list)

	/*
		监听所有执行中的任务
			给每一个任务单独开一个事件。
				监听某个具体执行的任务
				任务完成之后主销事件
					需要确定一下客户端监听未注册事件的情况。
	*/

	_db := db.New(context.Background())
	envs := []*models.CronSetting{}
	mu := sync.Mutex{} // 互斥锁
	// 10分钟更新一次环境信息
	go func() {
		list := []*models.CronSetting{}
		_db.Raw("SELECT id, name FROM `cron_setting` WHERE scene=? and status=?", models.SceneEnv, enum.StatusActive).Find(&list)
		mu.Lock()
		envs = list
		mu.Unlock()
		time.Sleep(time.Minute * 10)
	}()

	for {
		curTime := time.Now()
		list := dm.cron.Entries()
		execList := map[string][]*dtos.ExecQueueItem{}
		registerList := map[string][]*dtos.ExecQueueItem{}
		for _, v := range list {
			c, ok := v.Job.(*JobConfig)
			refType := "config"
			if !ok {
				c2, ok := v.Job.(*JobPipeline)
				if !ok {
					continue
				}
				c = c2.GetConf()
				refType = "pipeline"
			}
			// 区分环境
			if c.isRun == true {
				execList[c.conf.Env] = append(execList[c.conf.Env], &dtos.ExecQueueItem{
					RefId:    c.conf.Id,
					RefType:  refType,
					EntryId:  c.conf.EntryId,
					Name:     c.conf.Name,
					Duration: curTime.Sub(c.runTime).Seconds(),
				})
			}
			registerList[c.conf.Env] = append(registerList[c.conf.Env], &dtos.ExecQueueItem{
				RefId:   c.conf.Id,
				RefType: refType,
				EntryId: c.conf.EntryId,
				Name:    c.conf.Name,
			})
		}
		// 执行发送
		// 没有任务还是要发送空数据，不然客户端也不知道上一次的任务是否完成。
		for _, v := range envs {
			// 执行队列
			execItem := append(execList[""], execList[v.Name]...)
			if len(execItem) > 0 {
				b, _ := jsoniter.MarshalToString(execItem)
				sse.Serve().SendEventMessage(v.Name+".exec.queue", b)
			} else {
				sse.Serve().SendEventMessage(v.Name+".exec.queue", `[]`)
			}
			// 注册队列
			registerItem := append(registerList[""], registerList[v.Name]...)
			if len(registerItem) > 0 {
				b, _ := jsoniter.MarshalToString(registerItem)
				sse.Serve().SendEventMessage(v.Name+".register.queue", b)
			} else {
				sse.Serve().SendEventMessage(v.Name+".register.queue", `[]`)
			}
		}
		time.Sleep(time.Second * 5) // 延迟
	}
}

// 添加任务
func (dm *TaskService) AddConfig(conf *models.CronConfig) error {
	var id cron.EntryID
	var err error
	if conf == nil {
		return errors.New("未指定任务")
	}
	j := NewJobConfig(conf)
	if conf.Type == models.TypeOnce {
		id, err = dm.addOnce(conf.Spec, j)
	} else {
		id, err = dm.cron.AddJob(conf.Spec, j)
	}
	log.Printf("添加任务 %v -> %v", conf.Id, id)
	if err != nil {
		_, span := tracing.Tracer(conf.Env+"-cronin", trace.WithInstrumentationAttributes(
			attribute.String("driver", "mysql"),
			attribute.String("env", conf.Env),
		)).Start(context.Background(), "job-"+conf.GetProtocolName(), trace.WithAttributes(attribute.Int("ref_id", conf.Id)))
		defer span.End()
		span.SetStatus(tracing.StatusError, "任务启动失败")
		span.AddEvent("error", trace.WithAttributes(attribute.String("error.object", err.Error())))
		return err
	}
	conf.EntryId = int(id)
	return nil
}

func (dm *TaskService) AddPipeline(conf *models.CronPipeline) error {
	var id cron.EntryID
	var err error
	if conf == nil {
		return errors.New("未指定任务")
	}
	j := NewJobPipeline(conf)
	if conf.Type == models.TypeOnce {
		id, err = dm.addOnce(conf.Spec, j)
	} else {
		id, err = dm.cron.AddJob(conf.Spec, j)
	}
	if err != nil {
		_, span := tracing.Tracer(conf.Env+"-cronin", trace.WithInstrumentationAttributes(
			attribute.String("driver", "mysql"),
			attribute.String("env", conf.Env),
		)).Start(context.Background(), "job-pipeline", trace.WithAttributes(attribute.Int("ref_id", conf.Id)))
		defer span.End()
		span.SetStatus(tracing.StatusError, "任务启动失败")
		span.AddEvent("error", trace.WithAttributes(attribute.String("error.object", err.Error())))
		return err
	}
	conf.EntryId = int(id)
	return nil
}

// 添加单次任务
func (dm *TaskService) addOnce(spec string, j cron.Job) (cron.EntryID, error) {
	s, err := NewScheduleOnce(spec)
	if err != nil {
		return 0, err
	}
	return dm.cron.Schedule(s, j), nil
}

// 删除任务
func (dm *TaskService) Del(conf *models.CronConfig) {
	log.Printf("移除任务 %v -> %v", conf.Id, conf.EntryId)
	if conf.EntryId == 0 {
		return
	}
	dm.cron.Remove(cron.EntryID(conf.EntryId))
}

// 删除任务
func (dm *TaskService) DelPipeline(conf *models.CronPipeline) {
	if conf.EntryId == 0 { // 这里其实应该到任务队列取找执行id，mysql找是下册
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
		Rpc:     &pb.CronRpc{},
		Cmd:     &pb.CronCmd{},
		Sql:     &pb.CronSql{},
		Jenkins: &pb.CronJenkins{},
		Git:     &pb.CronGit{},
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
