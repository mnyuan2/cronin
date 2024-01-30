package biz

import (
	"bytes"
	"context"
	"cron/internal/basic/conv"
	"cron/internal/basic/db"
	"cron/internal/basic/enum"
	"cron/internal/basic/errs"
	"cron/internal/basic/grpcurl"
	"cron/internal/basic/tracing"
	"cron/internal/basic/util"
	"cron/internal/biz/dtos"
	"cron/internal/data"
	"cron/internal/models"
	"cron/internal/pb"
	"crypto/tls"
	"fmt"
	"github.com/axgle/mahonia"
	"github.com/jhump/protoreflect/grpcreflect"
	jsoniter "github.com/json-iterator/go"
	"github.com/robfig/cron/v3"
	"go.opentelemetry.io/otel/attribute"
	codes2 "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"io"
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
	//if m.execTime.Unix()-60*60*1 < time.Now().Unix() {
	//	return nil, errs.New(nil, "执行时间至少间隔1小时以后")
	//}

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
	tracer       trace.Tracer
}

// 任务执行器
func NewCronJob(conf *models.CronConfig) *CronJob {
	job := &CronJob{
		conf:         conf,
		commandParse: &pb.CronConfigCommand{},
		msgSetParse:  &dtos.MsgSetParse{MsgIds: []int{}, StatusIn: map[int]any{}, NotifyUserIds: []int{}, Set: []*pb.CronMsgSet{}},
	}

	_ = jsoniter.Unmarshal(conf.Command, job.commandParse)
	_ = jsoniter.Unmarshal(conf.MsgSet, &job.msgSetParse.Set)

	for _, s := range job.msgSetParse.Set {
		if s.Status == 0 || s.Status == enum.StatusDisable {
			job.msgSetParse.StatusIn[enum.StatusDisable] = struct{}{}
		}
		if s.Status == 0 || s.Status == enum.StatusActive {
			job.msgSetParse.StatusIn[enum.StatusActive] = struct{}{}
		}
		job.msgSetParse.NotifyUserIds = append(job.msgSetParse.NotifyUserIds, s.NotifyUserIds...)
		job.msgSetParse.MsgIds = append(job.msgSetParse.MsgIds, s.MsgId)
	}
	// 日志
	job.tracer = tracing.Tracer(job.conf.Env+"-cronin", trace.WithInstrumentationAttributes(
		attribute.String("driver", "mysql"),
		attribute.String("env", job.conf.Env),
		attribute.Int64("nonce", int64(job.conf.Id)),
	))

	return job
}

// 执行任务
func (job *CronJob) Run() {
	st := time.Now()
	ctx, span := job.tracer.Start(context.Background(), "job-"+job.conf.GetProtocolName())
	defer func() {
		if err := util.PanicInfo(recover()); err != "" {
			span.SetStatus(codes2.Error, "执行异常")
			span.AddEvent("x", trace.WithAttributes(attribute.String("error.object", err)))
		}
		span.End()
	}()
	span.SetAttributes(
		attribute.String("env", job.conf.Env),
		attribute.Int("config_id", job.conf.Id),
		attribute.String("protocol_name", job.conf.GetProtocolName()),
		attribute.String("component", "job"),
	)

	if job.conf.Type == models.TypeOnce { // 单次任务，自行移除
		e := cronRun.Entry(cron.EntryID(job.conf.EntryId))
		cronRun.Remove(cron.EntryID(job.conf.EntryId))
		if e.ID == 0 {
			span.SetStatus(codes2.Error, "重复执行？")
			return
		}
	}

	res, err := job.Exec(ctx)
	if err != nil {
		job.ErrorCount++
		job.messagePush(ctx, enum.StatusDisable, err.Error(), res, time.Since(st).Seconds())
	} else {
		job.ErrorCount = 0
		span.AddEvent("x", trace.WithAttributes(attribute.String("resp", string(res))))
		job.messagePush(ctx, enum.StatusActive, "ok", res, time.Since(st).Seconds())
	}

	// 连续错误达到5次，任务终止。
	e, ok := err.(*errs.Error)
	if job.ErrorCount >= 5 || (ok && e.Code() == errs.SysError.String()) {
		cronRun.Remove(cron.EntryID(job.conf.EntryId))
		job.conf.Status = models.ConfigStatusError
		job.conf.EntryId = 0
		if err := data.NewCronConfigData(ctx).ChangeStatus(job.conf, "执行失败"); err != nil {
			span.AddEvent("x", trace.WithAttributes(attribute.String("error", "任务状态写入失败"+err.Error())))
		}
	} else if job.conf.Type == models.TypeOnce { // 单次执行完毕后，状态也要更新
		cronRun.Remove(cron.EntryID(job.conf.EntryId))
		job.conf.Status = models.ConfigStatusFinish
		job.conf.EntryId = 0
		if err := data.NewCronConfigData(ctx).ChangeStatus(job.conf, "执行完成"); err != nil {
			span.AddEvent("x", trace.WithAttributes(attribute.String("error", "完成状态写入失败"+err.Error())))
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
		err = job.sqlFunc(ctx)
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
		return nil, errs.New(nil, "http method is empty", errs.SysError)
	}
	return job.httpRequest(ctx, method, http.Url, []byte(http.Body), header)
}

// rpc 执行函数
func (job *CronJob) rpcFunc(ctx context.Context) (res []byte, err error) {
	switch job.commandParse.Rpc.Method {
	case "GRPC":
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
		_, span := job.tracer.Start(ctx, "cmd-cmd")
		defer span.End()
		span.SetAttributes(
			attribute.String("component", "cmd"),
		)
		span.AddEvent("x", trace.WithAttributes(attribute.String("statement", job.commandParse.Cmd)))

		data := strings.Split(job.commandParse.Cmd, " ")
		if len(data) < 2 {
			return nil, errs.New(nil, "命令参数不合法，已跳过")
		}

		cmd := exec.Command(data[0], data[1:]...) // 合并 winds 命令
		if res, err = cmd.Output(); err != nil {
			span.SetStatus(codes2.Error, err.Error())
			return nil, err
		} else {
			srcCoder := mahonia.NewDecoder("gbk").ConvertString(string(res))
			span.AddEvent("x", trace.WithAttributes(attribute.String("console", srcCoder)))
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
		_, span := job.tracer.Start(ctx, "cmd-bash")
		defer span.End()
		span.SetAttributes(
			attribute.String("component", "cmd"),
		)
		span.AddEvent("x", trace.WithAttributes(attribute.String("statement", job.commandParse.Cmd)))

		cmd := exec.Command("/bin/bash", "-c", job.commandParse.Cmd)
		if res, err = cmd.Output(); err != nil {
			span.SetStatus(codes2.Error, err.Error())
			return nil, errs.New(err, "执行结果错误")
		} else {
			span.AddEvent("x", trace.WithAttributes(attribute.String("console", string(res))))
			return res, nil
		}
	}
}

// rpc 执行函数
func (job *CronJob) sqlFunc(ctx context.Context) (err error) {
	switch job.commandParse.Sql.Driver {
	case models.SqlSourceMysql:
		return job.sqlMysql(ctx, job.commandParse.Sql)
	default:
		return errs.New(nil, fmt.Sprintf("未支持的sql 驱动 %s", job.commandParse.Sql.Driver), errs.SysError)
	}
}

// http请求
func (job *CronJob) httpRequest(ctx context.Context, method, url string, body []byte, header map[string]string) (resp []byte, err error) {
	ctx, span := job.tracer.Start(ctx, "http-request")
	defer span.End()
	span.SetAttributes(
		attribute.String("component", "HTTP"),
		attribute.String("method", method),
	)

	h, _ := jsoniter.Marshal(header)
	span.AddEvent("x", trace.WithAttributes(
		attribute.String("req_header", string(h)),
		attribute.String("request", string(body)),
	))

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, errs.New(err, "请求构建失败")
	}
	for k, v := range header {
		req.Header.Set(k, v)
	}

	// 创建 HTTP 客户端
	client := &http.Client{
		Transport: &http.Transport{
			//MaxIdleConns:    10,
			//MaxConnsPerHost: 10,
			//IdleConnTimeout: 10 * time.Second,
			TLSClientConfig: &tls.Config{
				// 指定不校验 SSL/TLS 证书
				InsecureSkipVerify: true,
			},
		},
		//Timeout:   15 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, errs.New(err, "请求执行失败")
	}
	defer res.Body.Close()

	resp, err = io.ReadAll(res.Body)
	if err == nil && res.StatusCode != http.StatusOK {
		err = errs.New(fmt.Errorf("%v %s", res.StatusCode, http.StatusText(res.StatusCode)), "响应错误")
		span.SetStatus(codes2.Error, err.Error())
	}

	h, _ = jsoniter.Marshal(res.Header)
	span.AddEvent("x", trace.WithAttributes(
		attribute.String("resp_header", string(h)),
		attribute.String("response", string(resp)),
	))
	return resp, err
}

// grpc调用
func (job *CronJob) rpcGrpc(ctx context.Context, r *pb.CronRpc) (resp []byte, err error) {
	ctx, span := job.tracer.Start(ctx, "rpc-grpc")
	span.SetAttributes(
		attribute.String("component", "grpc-client"),
		attribute.String("target", r.Addr),
	)
	defer func() {
		if err != nil {
			span.SetStatus(codes2.Error, err.Error())
		}
		span.End()
	}()

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
		span.SetStatus(codes2.Error, err.Error())
	}

	span.SetAttributes(attribute.String("method", h.Method))
	span.AddEvent("x", trace.WithAttributes(
		attribute.String("req_header", string(h.GetSendHeadersMarshal())),
		attribute.String("resp_header", string(h.GetReceiveHeadersMarshal())),
		attribute.String("req", string(h.ReqMessages)),
		attribute.String("resp", h.RespMessages)),
	)
	return []byte(h.RespMessages), err
}

// 发送消息
func (job *CronJob) messagePush(ctx context.Context, status int, statusDesc string, body []byte, duration float64) {
	if _, ok := job.msgSetParse.StatusIn[status]; !ok {
		return
	}

	ctx, span := job.tracer.Start(ctx, "messagePush")
	defer func() {
		if err := util.PanicInfo(recover()); err != "" {
			span.SetStatus(codes2.Error, "执行异常")
			span.AddEvent("x", trace.WithAttributes(attribute.String("error.object", err)))
		}
		span.End()
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
		"log.status_name":      models.LogStatusMap[status],
		"log.status_desc":      statusDesc,
		"log.body":             strings.ReplaceAll(string(body), `"`, `\"`), // 内部存在双引号会引发错误
		"log.duration":         conv.Float64s().ToString(duration, 3),
		"log.create_dt":        time.Now().Format(time.DateTime),
		"user.username":        "",
		"user.mobile":          "",
	}

	for _, set := range job.msgSetParse.Set {
		if set.Status > 0 && set.Status != status {
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

		res, err := job.messagePushItem(ctx, []byte(msg.Content), args)
		if err != nil {
			span.SetStatus(codes2.Error, err.Error())
		} else {
			span.AddEvent("x", trace.WithAttributes(attribute.String("response", string(res))))
		}
	}
}

// 消息发送
func (job *CronJob) messagePushItem(ctx context.Context, templateByte []byte, args map[string]string) (res []byte, err error) {
	ctx, span := job.tracer.Start(ctx, "messagePushItem")
	defer span.End()
	// 变量替换
	for k, v := range args {
		templateByte = bytes.Replace(templateByte, []byte("[["+k+"]]"), []byte(v), -1)
	}

	template := &pb.SettingMessageTemplate{Http: &pb.CronHttp{}}
	if err = jsoniter.Unmarshal(templateByte, template); err != nil {
		return nil, errs.New(err, "消息模板解析错误")
	}

	// 执行推送
	return job.httpFunc(ctx, template.Http)
}
