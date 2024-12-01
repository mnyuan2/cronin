package biz

import (
	"bytes"
	"context"
	"cron/internal/basic/auth"
	"cron/internal/basic/conv"
	"cron/internal/basic/db"
	"cron/internal/basic/enum"
	"cron/internal/basic/errs"
	"cron/internal/basic/grpcurl"
	"cron/internal/basic/host"
	"cron/internal/basic/tracing"
	"cron/internal/basic/util"
	"cron/internal/biz/dtos"
	"cron/internal/data"
	"cron/internal/models"
	"cron/internal/pb"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/axgle/mahonia"
	"github.com/jhump/protoreflect/grpcreflect"
	jsoniter "github.com/json-iterator/go"
	"github.com/robfig/cron/v3"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
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
		return nil, errs.New(err, "执行时间格式不规范", errs.ParamError)
	}
	if m.execTime.Unix()-60 < time.Now().Unix() {
		return nil, errs.New(nil, "执行时间不得小于当前1分钟", errs.ParamError)
	}

	return m, nil
}

func NewScheduleOnceToTime(t time.Time) (m *ScheduleOnce, err error) {
	m = &ScheduleOnce{
		execTime: t,
	}

	if m.execTime.IsZero() {
		return nil, errs.New(err, "执行时间格式不规范", errs.ParamError)
	}
	return m, nil
}

// Next 返回下一次执行时间
//
//	如果当前时间大于执行时间返回空时间，将不会触发执行；否则将会按返回的时间执行。
func (m *ScheduleOnce) Next(t time.Time) time.Time {
	if m.execTime.UnixNano() < t.UnixNano() {
		return time.Time{}
	}
	return m.execTime
}

type JobConfig struct {
	conf         *models.CronConfig
	commandParse *pb.CronConfigCommand
	msgSetParse  *dtos.MsgSetParse
	ErrorCount   int // 连续错误
	tracer       trace.Tracer
	varParams    map[string]any
	isRun        bool // 是否执行中，false.否、true.是
	runTime      time.Time
	ctxCancel    context.CancelCauseFunc // 上下文取消
}

// 任务执行器
func NewJobConfig(conf *models.CronConfig) *JobConfig {
	job := &JobConfig{
		conf: conf,
	}

	// 日志
	job.tracer = tracing.Tracer(job.conf.Env+"-cronin", trace.WithInstrumentationAttributes(
		attribute.String("driver", "mysql"),
		attribute.String("env", job.conf.Env),
	))

	return job
}

// 解析参数
func (job *JobConfig) ParseParams(in map[string]any) (map[string]any, errs.Errs) {
	job.varParams = map[string]any{}
	if len(job.conf.VarFields) > 5 {
		// 参数也可以通过模板初始化，以获得动态默认值
		str, err := conv.DefaultStringTemplate().Execute(job.conf.VarFields)
		if err != nil {
			return nil, errs.New(err, "模板错误")
		}

		temp := []*pb.KvItem{}
		if err := jsoniter.Unmarshal(str, &temp); err != nil {
			return nil, errs.New(err, "变量参数字段解析错误")
		}
		for _, item := range temp {
			if item.Key == "" {
				continue
			}
			job.varParams[item.Key] = item.Value
		}
	}

	for k, v := range in {
		if _, ok := job.varParams[k]; ok {
			job.varParams[k] = v
		}
	}
	return job.varParams, nil
}

// 解析
func (job *JobConfig) Parse(params map[string]any) errs.Errs {
	// 进行模板替换
	cmd, err := conv.DefaultStringTemplate().SetParam(params).Execute(job.conf.Command)
	if err != nil {
		return errs.New(err, "模板错误")
	}

	job.commandParse = &pb.CronConfigCommand{}
	job.msgSetParse = &dtos.MsgSetParse{MsgIds: []int{}, StatusIn: map[int]any{}, NotifyUserIds: []int{}, Set: []*pb.CronMsgSet{}}
	if cmd != nil {
		if err := jsoniter.Unmarshal(cmd, job.commandParse); err != nil {
			return errs.New(err, "配置解析错误")
		}
	}
	if job.conf.MsgSet != nil {
		if err := jsoniter.Unmarshal(job.conf.MsgSet, &job.msgSetParse.Set); err != nil {
			return errs.New(err, "消息设置解析错误")
		}
	}

	for _, s := range job.msgSetParse.Set {
		job.msgSetParse.StatusIn[s.Status] = struct{}{}
		job.msgSetParse.NotifyUserIds = append(job.msgSetParse.NotifyUserIds, s.NotifyUserIds...)
		job.msgSetParse.MsgIds = append(job.msgSetParse.MsgIds, s.MsgId)
	}
	return nil
}

// 执行任务
func (job *JobConfig) Run() {
	var err errs.Errs
	var res []byte
	ctx1, span := job.tracer.Start(context.Background(), "job-"+job.conf.GetProtocolName(), trace.WithAttributes(attribute.Int("ref_id", job.conf.Id)))
	defer func() {
		job.isRun = false
		if res != nil {
			span.AddEvent("", trace.WithAttributes(attribute.String("resp", string(res))))
		}
		if err != nil {
			span.SetStatus(tracing.StatusError, err.Desc())
			span.AddEvent("error", trace.WithAttributes(attribute.String("error.object", err.Error())))
		} else {
			span.SetStatus(tracing.StatusOk, "")
		}
		if er := util.PanicInfo(recover()); er != "" {
			span.SetStatus(tracing.StatusError, "执行异常")
			span.AddEvent("error", trace.WithAttributes(attribute.String("error.panic", er)))
		}
		// 移除任务
		remark := ""
		g := data.NewChangeLogHandle(&auth.UserToken{UserName: "系统"}).SetType(models.LogTypeUpdateSys).OldConfig(*job.conf)
		if job.ErrorCount >= 10 {
			job.conf.Status = models.ConfigStatusError
			remark = "最大重试次数"
		} else if err != nil && err.Code() == errs.SysError.String() {
			job.conf.Status = models.ConfigStatusError
			remark = err.Desc() + "."
		} else if job.conf.Type == models.TypeOnce { // 单次执行完毕后，状态也要更新
			job.conf.Status = models.ConfigStatusFinish
			remark = "执行完成"
		}
		if remark != "" {
			rawId := cron.EntryID(job.conf.EntryId)
			e := cronRun.Entry(rawId)
			if e.ID == rawId {
				cronRun.Remove(e.ID)
			}
			job.conf.EntryId = 0
			if er := data.NewCronConfigData(ctx1).ChangeStatus(job.conf, remark); er != nil {
				span.AddEvent("error", trace.WithAttributes(attribute.String("error.object", "完成状态写入失败"+er.Error())))
			} else if er = data.NewCronChangeLogData(ctx1).Write(g.NewConfig(*job.conf)); er != nil {
				span.AddEvent("", trace.WithAttributes(attribute.String("warning", "变更记录写入失败"+er.Error())))
			}
		}
		span.End()
	}()
	if job.isRun {
		err = errs.New(nil, "任务正在进行中，跳过")
		return
	}
	ctx, cancel := context.WithCancelCause(ctx1)
	job.ctxCancel = cancel
	job.isRun = true
	job.runTime = time.Now()
	span.SetAttributes(
		attribute.String("env", job.conf.Env),
		attribute.String("config_name", job.conf.Name),
		attribute.String("protocol_name", job.conf.GetProtocolName()),
		attribute.String("component", "config"),
	)

	param, err := job.ParseParams(nil)
	if err = job.Parse(param); err != nil {
		return
	}

	for i := 0; i <= job.conf.ErrRetryNum; i++ {
		res, err = job.Exec(ctx)
		if err == nil {
			res, err = job.AfterTmpl(res, param)
		}

		if err != nil {
			job.ErrorCount++
			go job.messagePush(ctx, enum.StatusDisable, err.Desc(), []byte(err.Error()), time.Since(job.runTime).Seconds())
		} else {
			job.ErrorCount = 0
			go job.messagePush(ctx, enum.StatusActive, "ok", res, time.Since(job.runTime).Seconds())
			break // 成功后调出，不再重试
		}
	}

	if job.conf.AfterSleep > 0 {
		job.Sleep(ctx, time.Second*time.Duration(job.conf.AfterSleep))
	}
}

// 执行任务 未注册版本
//
//	给定 context 的情况下，内部不要开协程，会引发上下文终止的错误
func (job *JobConfig) Running(ctx context.Context, remark string, params map[string]any) (res []byte, err errs.Errs) {
	job.isRun = true
	job.runTime = time.Now()
	ctx, span := job.tracer.Start(ctx, "job-"+job.conf.GetProtocolName(), trace.WithAttributes(attribute.Int("ref_id", job.conf.Id)))
	defer func() {
		job.isRun = false
		if res != nil {
			span.AddEvent("", trace.WithAttributes(attribute.String("resp", string(res))))
		}
		if err != nil {
			span.SetStatus(tracing.StatusError, err.Desc())
			span.AddEvent("error", trace.WithAttributes(attribute.String("error.object", err.Error())))
		} else {
			span.SetStatus(tracing.StatusOk, "")
		}
		if er := util.PanicInfo(recover()); er != "" {
			span.SetStatus(tracing.StatusError, "执行异常")
			span.AddEvent("error", trace.WithAttributes(attribute.String("error.panic", er)))
		}
		p, _ := jsoniter.Marshal(job.varParams)
		span.AddEvent("", trace.WithAttributes(attribute.Stringer("var_params", bytes.NewBuffer(p))))
		span.End()
	}()
	log.Println("任务执行 开始", job.conf.Name)
	span.SetAttributes(
		attribute.String("env", job.conf.Env),
		attribute.String("config_name", job.conf.Name),
		attribute.String("protocol_name", job.conf.GetProtocolName()),
		attribute.String("component", "config"),
		attribute.String("remark", remark),
	)

	params, err = job.ParseParams(params)
	if err != nil {
		span.AddEvent("", trace.WithAttributes(attribute.Stringer("config", bytes.NewBuffer(job.conf.Command))))
		return
	}
	if err = job.Parse(params); err != nil {
		span.AddEvent("", trace.WithAttributes(attribute.Stringer("config", bytes.NewBuffer(job.conf.Command))))
		return
	}

	for i := 0; i <= job.conf.ErrRetryNum; i++ {
		res, err = job.Exec(ctx)
		if err == nil {
			res, err = job.AfterTmpl(res, params)
		}
		if err != nil {
			job.ErrorCount++
			job.messagePush(ctx, enum.StatusDisable, err.Desc(), []byte(err.Error()), time.Since(job.runTime).Seconds())
		} else {
			job.ErrorCount = 0
			job.messagePush(ctx, enum.StatusActive, "ok", res, time.Since(job.runTime).Seconds())
			break
		}
	}

	if job.conf.AfterSleep > 0 {
		job.Sleep(ctx, time.Second*time.Duration(job.conf.AfterSleep))
	}
	return res, err
}

// 停止任务，Run 执行时有效
func (job *JobConfig) Stop(ctx context.Context, remark string) error {
	if !job.isRun {
		return nil
	}
	job.ctxCancel(errors.New(remark))
	return nil
}

func (job *JobConfig) Exec(ctx context.Context) (res []byte, err errs.Errs) {
	switch job.conf.Protocol {
	case models.ProtocolHttp:
		res, err = job.httpFunc(ctx, job.commandParse.Http)
	case models.ProtocolRpc:
		res, err = job.rpcFunc(ctx)
	case models.ProtocolCmd:
		res, err = job.cmdFunc(ctx, job.commandParse.Cmd)
	case models.ProtocolSql:
		err = job.sqlFunc(ctx)
	case models.ProtocolJenkins:
		err = job.jenkins(ctx, job.commandParse.Jenkins)
	case models.ProtocolGit:
		res, err = job.gitFunc(ctx, job.commandParse.Git)
	default:
		err = errs.New(nil, fmt.Sprintf("未支持的protocol=%v", job.conf.Protocol))
	}
	return res, err
}

// 结果模板处理
func (job *JobConfig) AfterTmpl(result []byte, param map[string]any) (out []byte, err errs.Errs) {
	p := map[string]any{
		"result": string(result),
	}
	for k, v := range param {
		p[k] = v
	}
	if job.conf.AfterTmpl != "" {
		str, er := conv.DefaultStringTemplate().
			SetParam(p).
			Execute([]byte(job.conf.AfterTmpl))
		if er != nil {
			return nil, errs.New(er, "结果 模板错误")
		}
		out = str
	}
	return out, nil
}

// 等待
func (job *JobConfig) Sleep(ctx context.Context, duration time.Duration) errs.Errs {
	select {
	case <-ctx.Done():
		err := context.Cause(ctx)
		if err == nil {
			err = ctx.Err()
		}
		return errs.New(err, "上下文完成")
	case <-time.After(duration):
		return nil
	}
}

// http 执行函数
func (job *JobConfig) httpFunc(ctx context.Context, http *pb.CronHttp) (res []byte, err errs.Errs) {
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
func (job *JobConfig) rpcFunc(ctx context.Context) (res []byte, err errs.Errs) {
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
func (job *JobConfig) cmdFunc(ctx context.Context, r *pb.CronCmd) (res []byte, err errs.Errs) {
	ctx, span := job.tracer.Start(ctx, "cmd-"+r.Type)
	defer func() {
		if res != nil {
			span.AddEvent("", trace.WithAttributes(attribute.String("console", string(res))))
		}
		if err != nil {
			span.SetStatus(tracing.StatusError, err.Desc())
			span.AddEvent("error", trace.WithAttributes(attribute.String("error.object", err.Error())))
		}
		span.End()
	}()
	span.SetAttributes(attribute.String("component", r.Type))
	//确认数据来源
	statement := ""
	if r.Origin == "git" {
		files, err := job.getGitFile(ctx, r.Statement.Git)
		if err != nil {
			return nil, err
		}
		if len(files) < 1 { // 仅支持单文件
			return nil, errs.New(nil, "未配置有效文件")
		} else if len(files) > 1 {
			return nil, errs.New(nil, "仅支持单文件")
		}
		statement = string(files[0].Byte)
	} else {
		statement = r.Statement.Local
	}
	span.AddEvent("", trace.WithAttributes(attribute.String("statement", statement)))

	// 远程执行
	if r.Host != nil && r.Host.Id > 0 {
		span.SetAttributes(attribute.Int("host.id", r.Host.Id))
		s := &pb.SettingSource{}
		source, er := data.NewCronSettingData(ctx).GetSourceOne(job.conf.Env, r.Host.Id)
		if er != nil {
			return nil, errs.New(er, "连接配置异常")
		}
		if er = jsoniter.UnmarshalFromString(source.Content, s); er != nil {
			return nil, errs.New(er, "连接配置解析异常")
		}
		span.AddEvent("x",
			trace.WithAttributes(attribute.String("host.name", source.Name)),
			trace.WithAttributes(attribute.String("host.ip", s.Host.Ip)))

		return host.NewHost(&host.Config{
			Ip:     s.Host.Ip,
			Port:   s.Host.Port,
			User:   s.Host.User,
			Secret: s.Host.Secret,
		}).RemoteExec(r.Type + " " + statement)
	}

	// 本地执行
	switch r.Type {
	case "cmd":
		args := strings.Split(statement, " ")
		if len(args) < 2 {
			return nil, errs.New(nil, "命令参数不合法，已跳过")
		}
		cmd := exec.Command(args[0], args[1:]...) // 合并 winds 命令
		if re, er := cmd.Output(); err != nil {
			return re, errs.New(er, "执行错误")
		} else {
			srcCoder := mahonia.NewDecoder("gbk").ConvertString(string(re))
			res = []byte(srcCoder)
		}

	case "sh":
		e := exec.Command("sh", "-c", statement)
		cmd, er := e.Output()
		if er != nil {
			return nil, errs.New(er, "执行结果错误")
		}
		srcCoder := mahonia.NewDecoder("gbk").ConvertString(string(cmd))
		res = []byte(srcCoder)

	case "bash":
		cmd := exec.Command("/bin/bash", "-c", statement)
		re, er := cmd.Output()
		if er != nil {
			return nil, errs.New(er, "执行结果错误")
		}
		res = re
	}
	return res, nil
}

// http请求
func (job *JobConfig) httpRequest(ctx context.Context, method, url string, body []byte, header map[string]string) (resp []byte, err errs.Errs) {
	ctx, span := job.tracer.Start(ctx, "http-request")
	defer func() {
		if err != nil {
			span.SetStatus(tracing.StatusError, err.Desc())
			span.AddEvent("error", trace.WithAttributes(attribute.String("error.object", err.Error())))
		}
		span.End()
	}()
	span.SetAttributes(
		attribute.String("component", "HTTP"),
		attribute.String("method", method),
	)

	h, _ := jsoniter.Marshal(header)
	span.AddEvent("", trace.WithAttributes(
		attribute.String("url", url),
		attribute.String("request_header", string(h)),
		attribute.String("request", string(body)),
	))

	req, er := http.NewRequest(method, url, bytes.NewBuffer(body))
	if er != nil {
		return nil, errs.New(er, "请求构建失败")
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

	res, er := client.Do(req)
	if er != nil {
		return nil, errs.New(er, "请求执行失败")
	}
	defer res.Body.Close()

	resp, er = io.ReadAll(res.Body)
	if er == nil && res.StatusCode != http.StatusOK {
		err = errs.New(fmt.Errorf("%v %s", res.StatusCode, http.StatusText(res.StatusCode)), "响应错误")
	}

	h, _ = jsoniter.Marshal(res.Header)
	span.AddEvent("", trace.WithAttributes(
		attribute.String("response_header", string(h)),
		attribute.String("response", string(resp)),
	))
	return resp, err
}

// grpc调用
func (job *JobConfig) rpcGrpc(ctx context.Context, r *pb.CronRpc) (resp []byte, err errs.Errs) {
	ctx, span := job.tracer.Start(ctx, "rpc-grpc")
	span.SetAttributes(
		attribute.String("component", "grpc-client"),
		attribute.String("target", r.Addr),
		attribute.String("action", r.Action),
	)
	defer func() {
		if err != nil {
			span.SetStatus(tracing.StatusError, err.Desc())
			span.AddEvent("error", trace.WithAttributes(attribute.String("error.object", err.Error())))
		}
		span.End()
	}()

	cli, er := grpcurl.BlockingDial(ctx, "tcp", r.Addr, nil)
	if er != nil {
		return nil, errs.New(er, fmt.Sprintf("拨号目标主机 %s 失败", r.Addr))
	}
	// 解析描述文件
	var descSource grpcurl.DescriptorSource
	if r.Proto != "" {
		fds, er := grpcurl.ParseProtoString(r.Proto)
		if er != nil {
			return nil, errs.New(er, "无法解析给定的proto文件")
		}
		descSource, er = grpcurl.DescriptorSourceFromFileDescriptors(fds...)
		if er != nil {
			return nil, errs.New(er, "proto描述解析错误")
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
	rf, formatter, er := grpcurl.RequestParserAndFormatter(grpcurl.Format("json"), descSource, in, grpcurl.FormatOptions{
		EmitJSONDefaultFields: true, // 是否json格式
		IncludeTextSeparator:  false,
		AllowUnknownFields:    true,
	})
	if er != nil {
		return nil, errs.New(er, "请求解析器错误")
	}

	h := grpcurl.NewMyEventHandler(formatter)
	// 发起请求
	er = grpcurl.InvokeRPC(ctx, descSource, cli, r.Action, r.Header, h, rf.Next)
	// 处理错误
	if er != nil {
		errStatus, _ := status.FromError(er)
		h.SetStatus(errStatus)
	}
	if h.GetStatus().Code() != codes.OK {
		err = errs.New(fmt.Errorf("code:%v code_name:%v message:%v", int32(h.GetStatus().Code()), h.GetStatus().Code().String(), h.GetStatus().Message()), "响应错误")
	}

	span.AddEvent("data", trace.WithAttributes(
		attribute.String("method", h.Method),
		attribute.String("request_header", string(h.GetSendHeadersMarshal())),
		attribute.String("response_header", string(h.GetReceiveHeadersMarshal())),
		attribute.String("request", r.Body),
		attribute.String("response", h.RespMessages)),
	)
	return []byte(h.RespMessages), err
}

// 发送消息
func (job *JobConfig) messagePush(ctx context.Context, status int, statusDesc string, body []byte, duration float64) {
	if _, ok := job.msgSetParse.StatusIn[status]; !ok {
		return
	}

	ctx, span := job.tracer.Start(ctx, "message-push")
	defer func() {
		if err := util.PanicInfo(recover()); err != "" {
			span.SetStatus(tracing.StatusError, "执行异常")
			span.AddEvent("", trace.WithAttributes(attribute.String("error.panic", err)))
		}
		span.End()
	}()

	w := db.NewWhere().
		Eq("scene", models.SceneMsg).
		In("id", job.msgSetParse.MsgIds, db.RequiredOption()).
		Eq("status", enum.StatusActive)
	msgs, er := data.NewCronSettingData(ctx).Gets(w)
	if er != nil || len(msgs) == 0 {
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
		"log.body":             strings.ReplaceAll(string(body), `"`, `\\\"`), // 内部存在双引号会引发错误
		"log.duration":         conv.Float64s().ToString(duration, 3),
		"log.create_dt":        time.Now().Format(time.DateTime),
		"user.username":        "",
		"user.mobile":          "",
	}

	for _, set := range job.msgSetParse.Set {
		if set.Status != status {
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
		args["user.username"], _ = jsoniter.MarshalToString(username)
		args["user.mobile"], _ = jsoniter.MarshalToString(mobile)
		args["user.username"] = strings.ReplaceAll(args["user.username"], `"`, `\"`)
		args["user.mobile"] = strings.ReplaceAll(args["user.mobile"], `"`, `\"`)

		res, err := job.messagePushItem(ctx, []byte(msg.Content), args)
		if err != nil {
			span.SetStatus(tracing.StatusError, err.Desc())
			span.AddEvent("error", trace.WithAttributes(attribute.String("error.object", err.Error())))
		} else {
			span.AddEvent("", trace.WithAttributes(attribute.String("response", string(res))))
		}
	}
}

// 消息发送
func (job *JobConfig) messagePushItem(ctx context.Context, templateByte []byte, args map[string]string) (res []byte, err errs.Errs) {
	ctx, span := job.tracer.Start(ctx, "message-push-item")
	defer span.End()
	// 变量替换
	for k, v := range args {
		templateByte = bytes.Replace(templateByte, []byte("[["+k+"]]"), []byte(v), -1)
	}

	template := &pb.SettingMessageTemplate{Http: &pb.CronHttp{}}
	if er := jsoniter.Unmarshal(templateByte, template); err != nil {
		return nil, errs.New(er, "消息模板解析错误")
	}

	// 执行推送
	return job.httpFunc(ctx, template.Http)
}
