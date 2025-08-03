package biz

import (
	"context"
	"cron/internal/basic/auth"
	"cron/internal/basic/conv"
	"cron/internal/basic/db"
	"cron/internal/basic/errs"
	"cron/internal/basic/tracing"
	"cron/internal/data"
	"cron/internal/models"
	"cron/internal/pb"
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"strconv"
	"time"
)

// 组件与操作名称的对应关系
var componentToOperation = map[string]string{
	"config":   "job-task",
	"pipeline": "job-pipeline",
	"receive":  "job-receive",
}

type CronLogService struct {
	ctx  context.Context
	user *auth.UserToken
}

func NewCronLogService(ctx context.Context, user *auth.UserToken) *CronLogService {
	return &CronLogService{
		ctx:  ctx,
		user: user,
	}
}

// List 列表
func (dm *CronLogService) List(r *pb.CronLogListRequest) (resp *pb.CronLogListResponse, err error) {
	w, err := dm.listParseRequest(r)
	if err != nil {
		return nil, err
	}

	mod := data.NewCronLogSpanIndexV2Data(dm.ctx)

	total, list, err := mod.List(w, r.Page, r.Limit)
	if err != nil {
		return nil, errs.New(err, "查询失败")
	}
	if len(list) == 0 {
		return &pb.CronLogListResponse{List: []*pb.CronLogListItem{}, Page: &pb.Page{Page: r.Page, Size: r.Limit}}, nil
	}
	ids := make([]string, len(list))
	for i, item := range list {
		ids[i] = item.TraceId
	}

	traGro, err := mod.TraceIdGroup(db.NewWhere().Eq("env", r.Env).In("trace_id", ids))
	if err != nil {
		return nil, errs.New(err, "聚合失败")
	}
	resp = &pb.CronLogListResponse{
		List: make([]*pb.CronLogListItem, len(list)),
		Page: &pb.Page{Total: total, Page: r.Page, Size: r.Limit},
	}
	for i, item := range list {
		resp.List[i] = dm.listToOut(item, traGro[item.TraceId])
	}

	return resp, err
}

// 解析列表查询
func (dm *CronLogService) listParseRequest(r *pb.CronLogListRequest) (resp *db.Where, err error) {
	tags := map[string]any{}
	if r.Tags != "" {
		if err := jsoniter.UnmarshalFromString(r.Tags, &tags); err != nil {
			return nil, errs.New(err, "tags传递不规范")
		}
	}

	indexWhere := db.NewWhere()
	//where := db.NewWhere()
	//if r.Env != "" {
	indexWhere.Eq("env", r.Env)
	//where.Eq("env", r.Env)
	//} else {
	//	indexWhere.In("env", []string{dm.user.Env, ""})
	//where.In("env", []string{dm.user.Env, ""})
	//}
	if r.Page <= 0 {
		r.Page = 1
	}
	if r.Limit <= 15 {
		r.Limit = 15
	}

	for k, v := range tags {
		if k == "ref_id" {
			indexWhere.Eq("ref_id", conv.NewStr().ToString(v))
			continue
		} else if k == "component" {
			if op, ok := componentToOperation[v.(string)]; ok {
				indexWhere.Eq("operation", op)
				//where.Eq("operation", op)
				continue
			}
		}
		//where.Like("tags_kv", fmt.Sprintf("%s=%v", k, v))
	}
	indexWhere.Eq("operation", r.Operation)
	if r.TimestampStart != "" {
		startTime, err := time.ParseInLocation(time.DateTime, r.TimestampStart, time.Local)
		if err != nil {
			return nil, errors.New("开始时间不规范，请输入 yyyy-MM-dd hh:mm:ss 格式")
		}
		indexWhere.Gte("timestamp", startTime.UnixMicro())
	}
	if r.TimestampEnd != "" {
		endTime, err := time.ParseInLocation(time.DateTime, r.TimestampEnd, time.Local)
		if err != nil {
			return nil, errors.New("开始时间不规范，请输入 yyyy-MM-dd hh:mm:ss 格式")
		}
		indexWhere.Lte("timestamp", endTime.UnixMicro())
	}
	if r.DurationStart != "" {
		durationStart, err := strconv.ParseFloat(r.DurationStart, 64)
		if err != nil {
			return nil, errors.New("最小耗时错误，" + err.Error())
		}
		indexWhere.Gte("duration", (time.Duration(durationStart) * time.Second).Microseconds())
	}
	if r.DurationEnd != "" {
		durationEnd, err := strconv.ParseFloat(r.DurationEnd, 64)
		if err != nil {
			return nil, errors.New("最大耗时错误，" + err.Error())
		}
		indexWhere.Lte("duration", (time.Duration(durationEnd) * time.Second).Microseconds())
	}
	if r.Status != "" {
		status, err := strconv.Atoi(r.Status)
		if err != nil {
			return nil, errors.New("状态输入不规范，" + err.Error())
		}
		if status == int(tracing.StatusOk) {
			indexWhere.In("status", []int{int(tracing.StatusOk), int(tracing.StatusUnset)})
		} else if status == int(tracing.StatusError) {
			indexWhere.Eq("status", tracing.StatusError)
		}
	}
	if r.RefId != 0 {
		indexWhere.Eq("ref_id", conv.NewStr().ToString(r.RefId))
	}
	return indexWhere, nil
}

// Trace 踪迹
func (dm *CronLogService) Trace(r *pb.CronLogTraceRequest) (resp *pb.CronLogTraceResponse, err error) {
	if r.TraceId == "" {
		return nil, errs.New(nil, "未指定traceId")
	}

	w := db.NewWhere().In("env", []string{dm.user.Env, ""}).Eq("trace_id", r.TraceId)
	list, err := data.NewCronLogSpanData(dm.ctx).List(w, 1000, "*")

	// 树
	resp = &pb.CronLogTraceResponse{
		List:  []*pb.CronLogTraceItem{},
		Limit: 1000,
		Total: len(list),
	}

	tra := &pb.CronLogTraceItem{
		TraceId: r.TraceId,
		Spans:   []*pb.CronLogSpan{},
	}
	for _, item := range list {
		span := dm.toOut(item)
		tra.Spans = append(tra.Spans, span)
	}
	resp.List = append(resp.List, tra)

	return resp, err
}

// Del 删除
func (dm *CronLogService) Del(r *pb.CronLogDelRequest) (resp *pb.CronLogDelResponse, err error) {
	if r.Retention == "" {
		return nil, fmt.Errorf("retention 参数为必须")
	}

	re, err := time.ParseDuration(r.Retention)
	if err != nil {
		return nil, fmt.Errorf("retention 参数有误, %s", err.Error())
	} else if re.Hours() < 24 {
		return nil, fmt.Errorf("retention 参数不得小于24h")
	}
	end := time.Now().Add(-re)
	resp = &pb.CronLogDelResponse{}
	w := db.NewWhere().Lte("timestamp", end.UnixMicro())
	resp.Count, err = data.NewCronLogSpanData(dm.ctx).Del(w)
	if resp.Count > 0 {
		data.NewCronLogSpanIndexV2Data(dm.ctx).Del(db.NewWhere().Lte("timestamp", end.UnixMicro()))
	}

	return resp, err
}

// 转输出
func (dm *CronLogService) toOut(in *models.CronLogSpan) *pb.CronLogSpan {
	out := &pb.CronLogSpan{
		Timestamp:    in.Timestamp,
		Duration:     in.Duration,
		Status:       in.Status,
		StatusName:   models.LogSpanStatusMap[in.Status],
		StatusDesc:   in.StatusDesc,
		TraceId:      in.TraceId,
		SpanId:       in.SpanId,
		ParentSpanId: in.ParentSpanId,
		Service:      in.Service,
		Operation:    in.Operation,
		Tags:         []*pb.CronLogSpanKV{},
		Logs:         []*pb.CronLogSpanLog{},
	}

	jsoniter.Unmarshal(in.Tags, &out.Tags)
	jsoniter.Unmarshal(in.Logs, &out.Logs)

	return out
}

// 转输出
func (dm *CronLogService) listToOut(in *models.CronLogSpanIndexV2, gro []*data.TraceIdGroupItem) *pb.CronLogListItem {
	ti := time.UnixMicro(in.Timestamp)
	out := &pb.CronLogListItem{
		Timestamp:  ti.Format(time.DateTime),
		Duration:   in.Duration,
		Status:     in.Status,
		StatusName: models.LogSpanStatusMap[in.Status],
		//StatusDesc:   in.StatusDesc,
		TraceId: in.TraceId,
		SpanId:  in.SpanId,
		//Service:      in.Service,
		Operation: in.Operation,
		RefName:   in.RefName,
		RefId:     in.RefId,
	}
	if len(gro) > 0 {
		out.SpanTotal = 0
		out.RefId = gro[0].RefId
		out.RefName = gro[0].RefName
		out.Operation = gro[0].Operation
		for _, item := range gro {
			out.SpanGroup = append(out.SpanGroup, &pb.KvItem{
				Key:   item.Operation,
				Value: strconv.Itoa(item.TotalNum),
			})
			out.SpanTotal += int64(item.TotalNum)
		}
	}

	return out
}
