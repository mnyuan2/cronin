package tracing

// 节点配置
type SpanOption interface {
	parse() interface{}
}

//// 节点提取父及引用配置
//type SpanExtractOption string
//
//// 解析为内部的配置项
//func (o SpanExtractOption) parse() interface{} {
//	pCtx, err := opentracing.GlobalTracer().Extract(opentracing.TextMap, opentracing.TextMapCarrier{TraceId: string(o)})
//	if err != nil {
//		return nil
//	}
//	return ext.RPCServerOption(pCtx)
//}

//// http链路头提取
//// http默认大写开头，traceid默认小写需要转换兼容一下
//type HttpExtractOption http.Header
//
//func (o HttpExtractOption) parse() interface{} {
//	for k, v := range o {
//		if strings.ToLower(k) == TraceId {
//			o[TraceId] = v
//			break
//		}
//	}
//	if _, ok := o[TraceId]; !ok {
//		return nil
//	}
//	ctx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(o))
//	if err != nil {
//		return nil
//	}
//	return ext.RPCServerOption(ctx)
//}

//// kafka链路头提取
//type KafkaSpanExtractOption []kafka.Header
//
//func (o KafkaSpanExtractOption) parse() interface{} {
//	ctx, err := opentracing.GlobalTracer().Extract(opentracing.TextMap, o)
//	if err != nil {
//		return nil
//	}
//	return ext.RPCServerOption(ctx)
//}
//
//// 实现接口(opentracing.TextMapReader);
//func (o KafkaSpanExtractOption) ForeachKey(handler func(key, val string) error) error {
//	for _, val := range o {
//		if err := handler(val.Key, string(val.Value)); err != nil {
//			return err
//		}
//	}
//	return nil
//}
//
//// 实现接口(opentracing.TextMapWriter);
//func (o KafkaSpanExtractOption) Set(key, val string) {
//	o = append(o, kafka.Header{
//		Key:   key,
//		Value: []byte(val),
//	})
//}

//// 自定义链路开始时间
//type StartTimeOption time.Time
//
//func (o StartTimeOption) parse() interface{} {
//	return opentracing.StartTime(o)
//}
