package tracing

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"sync"
)

var (
	// 声明一个sync.Mutex 类型

	// 日 节点计数器
	dayCount = struct {
		mu         sync.Mutex
		Day        int
		TraceCount int64
		SpanCount  int64
	}{
		Day: 0, TraceCount: 0, SpanCount: 0,
	}
)

// 工厂模式
// 假设我只支持内部驱动，先完成再完善
func Tracer(name string, opts ...trace.TracerOption) trace.Tracer {
	conf := trace.NewTracerConfig(opts...)
	set := conf.InstrumentationAttributes()
	if val, ok := set.Value("driver"); ok {
		if val.AsString() == "mysql" {
			tra := &MysqlTracer{
				service: name,
			}
			if env, ok := set.Value("env"); ok {
				tra.env = env.AsString()
			}
			if env, ok := set.Value("nonce"); ok {
				tra.nonce = env.AsInt64()
			}

			return tra
		}
	}

	return otel.Tracer(name, opts...)
}

// 非全局性日志
func SqlTracer(name string, opts ...trace.TracerOption) *MysqlTracer {
	conf := trace.NewTracerConfig(opts...)
	set := conf.InstrumentationAttributes()

	tra := &MysqlTracer{
		service: name,
		spans:   map[string][]*MysqlSpan{},
	}
	if env, ok := set.Value("env"); ok {
		tra.env = env.AsString()
	}
	if env, ok := set.Value("nonce"); ok {
		tra.nonce = env.AsInt64()
	}

	return tra
}
