package biz

import (
	"bytes"
	"context"
	"cron/internal/basic/conv"
	"cron/internal/basic/db"
	"cron/internal/basic/enum"
	"cron/internal/basic/errs"
	"cron/internal/basic/grpcurl"
	"cron/internal/basic/util"
	"cron/internal/biz/dtos"
	"cron/internal/data"
	"cron/internal/models"
	"cron/internal/pb"
	"errors"
	"fmt"
	"github.com/axgle/mahonia"
	"github.com/jhump/protoreflect/grpcreflect"
	jsoniter "github.com/json-iterator/go"
	"github.com/robfig/cron/v3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type ScheduleOnce struct {
	s        string // 单次时间标志 YYYY-mm-dd hh:ii:ss
	execTime time.Time
}

func NewScheduleOnce(dateTime string) (m *ScheduleOnce, err error) {
	m = &ScheduleOnce{
		s: dateTime,
	}
	m.execTime, err = time.ParseInLocation(time.DateTime, dateTime, time.Local)
	if err != nil {
		return nil, errs.New(err, "执行时间格式不规范")
	}
	if m.execTime.Unix()-60*60*1 < time.Now().Unix() {
		return nil, errs.New(nil, "执行时间至少间隔1小时以后")
	}

	return m, nil
}

func (m *ScheduleOnce) Next(t time.Time) time.Time {
	if m.execTime.Unix() < t.Unix() {
		return t
	}
	return m.execTime
}

type CronJob struct {
	conf         *models.CronConfig
	commandParse *pb.CronConfigCommand
	msgSetParse  *dtos.MsgSetParse
	ErrorCount   int // 连续错误
}

// 任务执行器
func NewCronJob(conf *models.CronConfig) *CronJob {
	job := &CronJob{
		conf:         conf,
		commandParse: &pb.CronConfigCommand{},
		msgSetParse:  &dtos.MsgSetParse{MsgIds: []int{}, NotifyUserIds: []int{}, Set: []*pb.CronMsgSet{}},
	}

	_ = jsoniter.Unmarshal(conf.Command, job.commandParse)
	_ = jsoniter.Unmarshal(conf.MsgSet, &job.msgSetParse.Set)

	for _, s := range job.msgSetParse.Set {
		job.msgSetParse.NotifyUserIds = append(job.msgSetParse.NotifyUserIds, s.NotifyUserIds...)
		job.msgSetParse.MsgIds = append(job.msgSetParse.MsgIds, s.MsgId)
	}

	return job
}

// 执行任务
func (job *CronJob) Run() {
	var g *models.CronLog
	st := time.Now()
	ctx := context.Background()
	if job.conf.Type == models.TypeOnce { // 单次任务，自行移除
		e := cronRun.Entry(cron.EntryID(job.conf.EntryId))
		cronRun.Remove(cron.EntryID(job.conf.EntryId))
		if e.ID == 0 {
			return
		}
	}

	defer func() {
		if err := util.PanicInfo(recover()); err != "" {
			g2 := models.NewErrorCronLog(job.conf, "", errs.New(errors.New(err), "执行异常"), st)
			if g == nil {
				g = g2
			} else {
				g.Status = g2.Status
				g.StatusDesc += " " + g2.StatusDesc
				g.Body += "\n-----------------------\n" + g2.Body
			}
		}
		// 执行告警推送
		job.messagePush(ctx, g)

		data.NewCronLogData(ctx).Add(g)
		if job.conf.Type == models.TypeOnce { // 单次执行完毕后，状态也要更新
			job.conf.Status = models.ConfigStatusFinish
			job.conf.EntryId = 0
			data.NewCronConfigData(ctx).ChangeStatus(job.conf, "执行完成")
		}
	}()

	//fmt.Println("执行 "+job.conf.GetProtocolName()+" 任务", job.conf.Id, job.conf.Name)
	res, err := job.Exec(ctx)

	if err != nil {
		g = models.NewErrorCronLog(job.conf, string(res), err, st)
		job.ErrorCount++
	} else {
		g = models.NewSuccessCronLog(job.conf, string(res), st)
		job.ErrorCount = 0
	}
	// 连续错误达到5次，任务终止。
	e, ok := err.(*errs.Error)
	if job.ErrorCount >= 5 || (ok && e.Code() == errs.SysError.String()) || job.conf.Type == models.TypeOnce {
		cronRun.Remove(cron.EntryID(job.conf.EntryId))
		job.conf.Status = models.ConfigStatusError
		job.conf.EntryId = 0
		err := data.NewCronConfigData(ctx).ChangeStatus(job.conf, "执行失败")
		if err != nil {
			log.Println("任务状态写入失败 id", job.conf.Id, err.Error())
		}
	}
}

func (job *CronJob) Exec(ctx context.Context) (res []byte, err error) {
	switch job.conf.Protocol {
	case models.ProtocolHttp:
		res, err = job.httpFunc(ctx, job.commandParse.Http)
	case models.ProtocolRpc:
		res, err = job.rpcFunc(ctx)
	case models.ProtocolCmd:
		res, err = job.cmdFunc(ctx)
	case models.ProtocolSql:
		res, err = job.sqlFunc(ctx)
	default:
		err = errs.New(nil, fmt.Sprintf("未支持的protocol=%v", job.conf.Protocol))
	}
	return res, err
}

// http 执行函数
func (job *CronJob) httpFunc(ctx context.Context, http *pb.CronHttp) (res []byte, err error) {
	header := map[string]string{}
	for _, head := range http.Header {
		if head.Key == "" {
			continue
		}
		header[head.Key] = head.Value
	}
	method := models.ProtocolHttpMethodMap()[http.Method]
	if method == "" {
		// 任务设置有问题，提出执行队列，记录日志。
		//job.ErrorCount = -2
		return nil, errs.New(nil, "http method is empty", errs.SysError)
	}
	return job.httpRequest(ctx, method, http.Url, []byte(http.Body), header)
}

// rpc 执行函数
func (job *CronJob) rpcFunc(ctx context.Context) (res []byte, err error) {
	switch job.commandParse.Rpc.Method {
	case "GRPC":
		// 进行grpc处理
		// 目前还存在问题，无法通用性的提交和接收参数！
		return job.rpcGrpc(ctx, job.commandParse.Rpc)
	case "RPC":
		return nil, errs.New(nil, fmt.Sprintf("未支持的rpc method，任务已终止。"), errs.SysError)
		// 手头目前没有rpc的服务，不好测试验证。
	default:
		return nil, errs.New(nil, fmt.Sprintf("未支持的rpc method %s，任务已终止。", job.commandParse.Rpc.Method), errs.SysError)
	}
}

// rpc 执行函数
func (job *CronJob) cmdFunc(ctx context.Context) (res []byte, err error) {
	if runtime.GOOS == "windows" {
		data := strings.Split(job.commandParse.Cmd, " ")
		if len(data) < 2 {
			return nil, errs.New(nil, "命令参数不合法，已跳过")
		}

		cmd := exec.Command(data[0], data[1:]...) // 合并 winds 命令
		if res, err = cmd.Output(); err != nil {
			return nil, err
		} else {
			srcCoder := mahonia.NewDecoder("gbk").ConvertString(string(res))
			return []byte(srcCoder), nil
		}

		// windows下安装了sh.exe 可以使用；就是git；但要添加环境变量 git/bin
		//cmd := exec.Command("sh.exe","-c", job.commandParse.Cmd)
		//if res, err = cmd.Output(); err != nil{
		//	return nil, err
		//}else {
		//	return res, nil
		//}

	} else {
		cmd := exec.Command("/bin/bash", "-c", job.commandParse.Cmd)
		if res, err = cmd.Output(); err != nil {
			return nil, errs.New(err, "执行结果错误")
		} else {
			return res, nil
		}
	}
}

// rpc 执行函数
func (job *CronJob) sqlFunc(ctx context.Context) (res []byte, err error) {
	switch job.commandParse.Sql.Driver {
	case models.SqlSourceMysql:
		return job.sqlMysql(ctx, job.commandParse.Sql)
	default:
		return nil, errs.New(nil, fmt.Sprintf("未支持的sql 驱动 %s", job.commandParse.Sql.Driver), errs.SysError)
	}
}

// http请求
func (job *CronJob) httpRequest(ctx context.Context, method, url string, body []byte, header map[string]string) (resp []byte, err error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, errs.New(err, "请求构建失败")
	}
	for k, v := range header {
		req.Header.Set(k, v)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errs.New(err, "请求执行失败")
	}
	defer res.Body.Close()

	resp, err = io.ReadAll(res.Body)
	if err == nil && res.StatusCode != http.StatusOK {
		err = errs.New(fmt.Errorf("%v %s", res.StatusCode, http.StatusText(res.StatusCode)), "响应错误")
	}
	return resp, err
}

// grpc调用
func (job *CronJob) rpcGrpc(ctx context.Context, r *pb.CronRpc) (resp []byte, err error) {
	cli, err := grpcurl.BlockingDial(ctx, "tcp", r.Addr, nil)
	if err != nil {
		return nil, errs.New(err, fmt.Sprintf("拨号目标主机 %s 失败", r.Addr))
	}
	// 解析描述文件
	var descSource grpcurl.DescriptorSource
	if r.Proto != "" {
		fds, err := grpcurl.ParseProtoString(r.Proto)
		if err != nil {
			return nil, errs.New(err, "无法解析给定的proto文件")
		}
		descSource, err = grpcurl.DescriptorSourceFromFileDescriptors(fds...)
		if err != nil {
			return nil, errs.New(err, "proto描述解析错误")
		}
	} else { // 大部分服务器是不支持服务端的反射解析的
		md := grpcurl.MetadataFromHeaders(r.Header)
		refCtx := metadata.NewOutgoingContext(ctx, md)
		refClient := grpcreflect.NewClientAuto(refCtx, cli)
		descSource = grpcurl.DescriptorSourceFromServer(ctx, refClient)
	}

	var in io.Reader
	if r.Body == "" {
		in = os.Stdin
	} else {
		in = strings.NewReader(r.Body)
	}

	// 如果不是详细输出，那么还可以在每个消息之间包含记录分隔符，这样输出就可以通过管道输送到另一个grpcurl进程
	// 请求参数处理方法，把原json参数根据描述文件进行了转义处理。
	rf, formatter, err := grpcurl.RequestParserAndFormatter(grpcurl.Format("json"), descSource, in, grpcurl.FormatOptions{
		EmitJSONDefaultFields: true, // 是否json格式
		IncludeTextSeparator:  false,
		AllowUnknownFields:    true,
	})
	if err != nil {
		return nil, errs.New(err, "请求解析器错误")
	}

	h := grpcurl.NewMyEventHandler(formatter)
	// 发起请求
	err = grpcurl.InvokeRPC(ctx, descSource, cli, r.Action, r.Header, h, rf.Next)
	// 处理错误
	if err != nil {
		errStatus, _ := status.FromError(err)
		h.SetStatus(errStatus)
	}
	if h.GetStatus().Code() != codes.OK {
		err = errs.New(fmt.Errorf("code:%v code_name:%v message:%v", int32(h.GetStatus().Code()), h.GetStatus().Code().String(), h.GetStatus().Message()), "响应错误")
	}

	return []byte(h.RespMessages), err
}

// 发送消息
func (job *CronJob) messagePush(ctx context.Context, g *models.CronLog) {
	defer func() {
		if err := util.PanicInfo(recover()); err != "" {
			g2 := models.NewErrorCronLog(job.conf, "", errs.New(errors.New(err), "执行异常"), time.Now())
			g.Status = g2.Status
			g.StatusDesc += " " + g2.StatusDesc
			g.Body += "\n------------\n" + g2.Body
		}
	}()

	w := db.NewWhere().
		Eq("scene", models.SceneMsg).
		In("id", job.msgSetParse.MsgIds, db.RequiredOption()).
		Eq("status", enum.StatusActive)
	msgs, err := data.NewCronSettingData(ctx).Gets(w)
	if err != nil || len(msgs) == 0 {
		return
	}
	msgMaps := map[int]*models.CronSetting{}
	for _, m := range msgs {
		msgMaps[m.Id] = m
	}

	users, _ := data.NewCronUserData(ctx).
		GetList(db.NewWhere().In("id", job.msgSetParse.NotifyUserIds, db.RequiredOption()))
	userMaps := map[int]*models.CronUser{}
	for _, user := range users {
		userMaps[user.Id] = user
	}

	// 重组临时变量，默认置空，有效的写入新值
	args := map[string]string{
		"env":                  job.conf.Env,
		"config.name":          job.conf.Name,
		"config.protocol_name": job.conf.GetProtocolName(),
		"log.status_name":      g.GetStatusName(),
		"log.status_desc":      g.StatusDesc,
		"log.body":             strings.ReplaceAll(g.Body, `"`, `\"`), // 内部存在双引号会引发错误
		"log.duration":         conv.Float64s().ToString(g.Duration),
		"log.create_dt":        g.CreateDt,
		"user.username":        "",
		"user.mobile":          "",
	}

	for _, set := range job.msgSetParse.Set {
		if set.Status > 0 && set.Status != g.Status {
			continue
		}

		// 查询模板
		msg, ok := msgMaps[set.MsgId]
		if !ok {
			continue // 消息模板不存在或未启用
		}

		username, mobile := []string{}, []string{}
		for _, userId := range set.NotifyUserIds {
			if user, ok := userMaps[userId]; ok {
				if user.Username != "" {
					username = append(username, user.Username)
				}
				if user.Mobile != "" {
					mobile = append(mobile, user.Mobile)
				}
			}
		}
		args["user.username"] = strings.Join(username, ",")
		args["user.mobile"] = strings.Join(mobile, ",")

		// 变量替换
		str := []byte(msg.Content)
		for k, v := range args {
			str = bytes.Replace(str, []byte("[["+k+"]]"), []byte(v), -1)
		}

		template := &pb.SettingMessageTemplate{Http: &pb.CronHttp{}}
		if err := jsoniter.Unmarshal(str, template); err != nil {
			g.StatusDesc += " 消息模板解析错误"
			g.Body += "\n------------\n" + err.Error()
			continue // 解析错误
		}

		// 执行推送
		res, err := job.httpFunc(ctx, template.Http)
		if err != nil {
			g.StatusDesc += " 推送失败"
			g.Body += "\n------------\n" + string(res)
		}
	}
}
