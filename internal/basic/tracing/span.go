package tracing

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Span struct {
	driver trace.Span
	ctx    context.Context
}

// 开始一个链路日志节点
func StartSpan(ctx context.Context, operationName string, opts ...SpanOption) *Span {
	s := &Span{
		ctx: ctx,
	}
	if parent := ctx.Value("mysql_span"); parent != nil {
		s.driver = parent.(*MysqlSpan)
	}

	return s
}

// 设置标签描述
func (s *Span) SetTags(data map[string]any) *Span {
	if s.driver != nil {
		list := []attribute.KeyValue{}
		for k, v := range data {
			switch v.(type) {
			case string:
				list = append(list, attribute.String(k, v.(string)))
			case int:
				list = append(list, attribute.Int(k, v.(int)))
			case int64:
				list = append(list, attribute.Int64(k, v.(int64)))
			case float64, float32:
				list = append(list, attribute.Float64(k, v.(float64)))
			}
		}
		s.driver.SetAttributes(list...)
	}
	return s
}

// 设置记录描述
func (s *Span) SetLogs(data map[string]any) *Span {
	if s.driver != nil {
		list := []attribute.KeyValue{}
		for k, v := range data {
			switch v.(type) {
			case string:
				list = append(list, attribute.String(k, v.(string)))
			case int:
				list = append(list, attribute.Int(k, v.(int)))
			case int64:
				list = append(list, attribute.Int64(k, v.(int64)))
			case float64, float32:
				list = append(list, attribute.Float64(k, v.(float64)))
			}
		}
		s.driver.AddEvent("event1", trace.WithAttributes(list...))
	}
	return s
}

// 链路结束
func (s *Span) End() {
	if s.driver != nil {
		s.driver.End()
	}
}
