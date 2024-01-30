package tracing

import (
	"context"
	"cron/internal/basic/db"
	"cron/internal/models"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/embedded"
	"log"
	"strconv"
	"strings"
	"time"
)

//const DriverMysql = attribute.String("driver", "mysql")

// 收集队列
var mysqlQueue chan models.CronLogSpan
var gen = &mysqlIDGenerator{}

// 链路日志收集写入程序
func MysqlCollectorListen() {
	mysqlQueue = make(chan models.CronLogSpan, 5000)
	exec := make(chan byte, 1)
	defer close(mysqlQueue)
	defer close(exec)
	go func() {
		for {
			time.Sleep(3 * time.Second)
			exec <- 1
		}
	}()

	// 延长3秒、或超过1000条写入。
	for {
		<-exec
		l := len(mysqlQueue)
		index := 1
		if l == 0 {
			continue
		} else if l > 1000 {
			l = 1000
			exec <- 1
		}

		list := []models.CronLogSpan{}
		for item := range mysqlQueue {
			list = append(list, item)
			if index >= l {
				break
			}
			index++
		}

		// 执行写入
		if err := db.New(context.Background()).CreateInBatches(list, 1000).Error; err != nil {
			log.Println("MysqlCollector 日志写入失败，", err.Error())
		}
	}

}

type mysqlTracer struct {
	embedded.Tracer

	service string
	env     string
	nonce   int64
}

func (t *mysqlTracer) tracer() {}

// 链路id生成
type mysqlIDGenerator struct {
	startTime time.Time
	env       string
	nonce     int64
}

func (t *mysqlIDGenerator) NewIDs(ctx context.Context) (trace.TraceID, trace.SpanID) {
	if t.nonce == 0 {
		day := t.startTime.Day()
		if dayCount.Day != day {
			dayCount.Day = day
			dayCount.TraceCount = 1
			dayCount.SpanCount = 1
		}
		t.nonce = dayCount.TraceCount
		dayCount.TraceCount++
	}

	id := fmt.Sprintf("%02.2s%010.10v%04.4v", t.env, t.startTime.Unix(), t.nonce)
	hex := fmt.Sprintf("%032x", id) // 32位
	traceID, _ := trace.TraceIDFromHex(hex)
	spanID := t.NewSpanID(ctx, traceID)
	return traceID, spanID
}

func (t *mysqlIDGenerator) NewSpanID(ctx context.Context, traceID trace.TraceID) trace.SpanID {
	day := t.startTime.Day()
	if dayCount.Day != day {
		dayCount.Day = day
		dayCount.TraceCount = 1
		dayCount.SpanCount = 1
	}
	nonce := dayCount.SpanCount
	dayCount.SpanCount++

	id := fmt.Sprintf("%02.2v%06.6v", dayCount.Day, nonce)
	spanIDHex := fmt.Sprintf("%016x", id) // 16位
	spanID, _ := trace.SpanIDFromHex(spanIDHex)
	return spanID
}

// ParseID 解析16进制字符串
func (t *mysqlIDGenerator) ParseID(hexStr string) (string, error) {
	// 将每个16进制字符转换为字节值，并转换为ASCII字符
	var normalStr strings.Builder
	for i := 0; i < len(hexStr); i += 2 {
		byteValue, err := strconv.ParseUint(hexStr[i:i+2], 16, 8) // 使用ParseUint代替ParseByte
		if err != nil {
			return "", fmt.Errorf("无效的16进制字符：%s", hexStr[i:i+2])
		}
		normalStr.WriteByte(byte(byteValue)) // 将字节值转换为ASCII字符并追加到normalStr中
	}

	return normalStr.String(), nil // 返回转换后的正常字符串
}

func (t *mysqlTracer) Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	conf := trace.NewSpanStartConfig(opts...)
	span := &MysqlSpan{
		service:   t.service,
		operation: spanName,
		env:       t.env,
		startTime: time.Now(),
		tags:      []attribute.KeyValue{},
		logs:      []trace.EventOption{},
	}
	span.tags = append(span.tags, conf.Attributes()...)

	gen := &mysqlIDGenerator{
		startTime: span.startTime,
		env:       t.env,
		nonce:     t.nonce,
	}
	if parent := ctx.Value("mysql_span"); parent != nil {
		span.traceId = parent.(*MysqlSpan).traceId
		span.spanId = gen.NewSpanID(ctx, span.traceId)
	} else {
		span.traceId, span.spanId = gen.NewIDs(ctx)
	}

	ctx = context.WithValue(ctx, "mysql_span", span)
	return ctx, span
}

// mysql 驱动的 Span节点
type MysqlSpan struct {
	embedded.Span
	sc trace.SpanContext

	traceId   trace.TraceID
	spanId    trace.SpanID
	service   string
	operation string
	env       string
	// startTime 开始时间
	startTime time.Time
	// endTime 结束时间
	endTime time.Time
	// status 状态
	status codes.Code
	// 标签集
	tags []attribute.KeyValue
	// 日志集
	logs []trace.EventOption
}

// SpanContext returns an empty span context.
func (s *MysqlSpan) SpanContext() trace.SpanContext { return s.sc }

// IsRecording always returns false.
func (*MysqlSpan) IsRecording() bool { return false }

// SetStatus does nothing.
func (s *MysqlSpan) SetStatus(status codes.Code, desc string) {
	s.status = status
}

func (s *MysqlSpan) SetLocalStatus(status int, desc string) {
	s.status = codes.Code(status)
}

// SetAttributes 设置标签
//
//	后续支持条件查询，单个key与val不得超过120个字符。
func (s *MysqlSpan) SetAttributes(kv ...attribute.KeyValue) {
	s.tags = append(s.tags, kv...)
}

// AddEvent 记录日志
//
//	不支持查询
func (s *MysqlSpan) AddEvent(name string, options ...trace.EventOption) {
	s.logs = append(s.logs, options...)
}

// End does nothing.
func (s *MysqlSpan) End(...trace.SpanEndOption) {
	s.endTime = time.Now()
	// 执行日志的写入
	data := &models.CronLogSpan{
		Timestamp: s.startTime.Format(time.DateTime),
		//TraceId:   s.traceId.String(),
		Service:   s.service,
		Operation: s.operation,
		Duration:  s.endTime.Sub(s.startTime).Seconds(),
		Status:    int32(s.status),
		Env:       s.env,
	}
	data.TraceId, _ = gen.ParseID(s.traceId.String())
	data.SpanId, _ = gen.ParseID(s.spanId.String())
	data.Tags, _ = jsoniter.Marshal(s.tags)
	data.Logs, _ = jsoniter.Marshal(s.logs)

	//log.Println(data)
	mysqlQueue <- *data
}

// RecordError does nothing.
func (*MysqlSpan) RecordError(error, ...trace.EventOption) {}

// SetName does nothing.
func (*MysqlSpan) SetName(string) {}

// TracerProvider returns a No-Op TracerProvider.
func (*MysqlSpan) TracerProvider() trace.TracerProvider {
	return nil
}
