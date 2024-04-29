package tracing

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"testing"
	"time"
)

func TestStartSpan(t *testing.T) {
	ctx := context.Background()

	// 服务(程序)名称
	tracer := otel.Tracer("服务名")

	ctx, span := tracer.Start(ctx, "parent", trace.WithAttributes())
	// 添加自定义属性/标记
	span.SetAttributes(attribute.Int("roll.value", 010))
	// 添加时间/日志
	span.AddEvent("event 1",
		trace.WithAttributes(attribute.String("key1", "value1")),
		trace.WithAttributes(attribute.String("key2", "value2")),
	)

	// 子集元素; 此处使用的parent节点的context
	ctx, span2 := tracer.Start(ctx, "child")

	//defer span.End()
	span.End()
	span2.End()

	t.Logf("end...")
}

func TestStartSpan2(t *testing.T) {
	go MysqlCollectorListen()

	ctx := context.Background()
	tracer := Tracer("public-cronin", trace.WithInstrumentationAttributes(
		attribute.String("driver", "mysql"),
		attribute.String("env", "public"),
	))

	ctx, span := tracer.Start(ctx, "parent", trace.WithAttributes(
		attribute.String("tag1", "value"),
	))
	// 添加自定义属性/标记
	span.SetAttributes(attribute.Int("tag2", 010))
	// 添加时间/日志
	span.AddEvent("event 1",
		trace.WithAttributes(attribute.String("log1", "value1")),
		trace.WithAttributes(attribute.String("log2", "value2")),
	)

	// 子集元素; 此处使用的parent节点的context
	ctx, span2 := tracer.Start(ctx, "child")

	//defer s.End()
	span.End()
	span2.End()

	time.Sleep(time.Second * 60)
	t.Logf("end...")
}

func TestDemo2(t *testing.T) {
	gen := &mysqlIDGenerator{}
	s, err := gen.ParseID("70753137303634333133393230303031")
	fmt.Printf("%s, %v\n", s, err)
	s, err = gen.ParseID("3030303030303031")
	fmt.Printf("%s, %v\n", s, err)
}

func TestOption(t *testing.T) {
	conf := trace.NewSpanStartConfig(trace.WithTimestamp(time.Now()))
	fmt.Println(conf.Timestamp(), conf.Timestamp().IsZero())

	conf2 := trace.NewSpanStartConfig()
	fmt.Println(conf2.Timestamp(), conf2.Timestamp().IsZero())
}
