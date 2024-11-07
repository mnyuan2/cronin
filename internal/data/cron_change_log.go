package data

import (
	"context"
	"cron/internal/basic/auth"
	"cron/internal/basic/db"
	"cron/internal/models"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"reflect"
	"strconv"
	"time"
)

// 日志过程
type ChangeLogHandle struct {
	user *auth.UserToken
	log  *models.CronChangeLog
	conf []*models.CronConfig
	line []*models.CronPipeline
	rece []*models.CronReceive
}

func NewChangeLogHandle(user *auth.UserToken) *ChangeLogHandle {
	return &ChangeLogHandle{
		user: user,
		log:  &models.CronChangeLog{},
		conf: make([]*models.CronConfig, 2),
		line: make([]*models.CronPipeline, 2),
		rece: make([]*models.CronReceive, 2),
	}
}
func (h *ChangeLogHandle) SetType(val int) *ChangeLogHandle {
	h.log.Type = val
	return h
}

func (h *ChangeLogHandle) OldConfig(m models.CronConfig) *ChangeLogHandle {
	if h.log.RefType != "" && h.log.RefType != "config" {
		panic("日志引用类型冲突")
	}
	h.log.RefType = "config"
	h.conf[0] = &m
	return h
}

func (h *ChangeLogHandle) NewConfig(m models.CronConfig) *ChangeLogHandle {
	h.conf[1] = &m
	return h
}

func (h *ChangeLogHandle) OldPipeline(m models.CronPipeline) *ChangeLogHandle {
	if h.log.RefType != "" && h.log.RefType != "pipeline" {
		panic("日志引用类型冲突")
	}
	h.log.RefType = "pipeline"
	h.line[0] = &m
	return h
}

func (h *ChangeLogHandle) NewPipeline(m models.CronPipeline) *ChangeLogHandle {
	h.line[1] = &m
	return h
}

func (h *ChangeLogHandle) OldReceive(m models.CronReceive) *ChangeLogHandle {
	if h.log.RefType != "" && h.log.RefType != "receive" {
		panic("日志引用类型冲突")
	}
	h.log.RefType = "receive"
	h.rece[0] = &m
	return h
}

func (h *ChangeLogHandle) NewReceive(m models.CronReceive) *ChangeLogHandle {
	h.rece[1] = &m
	return h
}

// 比较并输出差异
func (h *ChangeLogHandle) diffConfig(old, new *models.CronConfig) (content []*models.ChangeLogField) {
	content = []*models.ChangeLogField{}
	if old.EntryId != new.EntryId {
		content = append(content, &models.ChangeLogField{
			Field:      "entry_id",
			VType:      reflect.Int.String(),
			OldVal:     old.EntryId,
			NewVal:     new.EntryId,
			FieldName:  "执行编号",
			OldValName: strconv.Itoa(old.EntryId),
			NewValName: strconv.Itoa(new.EntryId),
		})
	}
	if old.Type != new.Type {
		content = append(content, &models.ChangeLogField{
			Field:      "type",
			VType:      reflect.Int.String(),
			OldVal:     old.Type,
			NewVal:     new.Type,
			FieldName:  "类型",
			OldValName: models.ConfigTypeMap[old.Type],
			NewValName: models.ConfigTypeMap[new.Type],
		})
	}
	if old.Name != new.Name {
		content = append(content, &models.ChangeLogField{
			Field:      "name",
			VType:      reflect.String.String(),
			OldVal:     old.Name,
			NewVal:     new.Name,
			FieldName:  "任务名称",
			OldValName: old.Name,
			NewValName: new.Name,
		})
	}
	if old.Spec != new.Spec {
		content = append(content, &models.ChangeLogField{
			Field:      "spec",
			VType:      reflect.String.String(),
			OldVal:     old.Spec,
			NewVal:     new.Spec,
			FieldName:  "执行时间",
			OldValName: old.Spec,
			NewValName: new.Spec,
		})
	}
	if old.Protocol != new.Protocol {
		content = append(content, &models.ChangeLogField{
			Field:      "protocol",
			VType:      reflect.Int.String(),
			OldVal:     old.Protocol,
			NewVal:     new.Protocol,
			FieldName:  "协议",
			OldValName: models.ProtocolMap[old.Protocol],
			NewValName: models.ProtocolMap[new.Protocol],
		})
	}
	if old.CommandHash != new.CommandHash {
		g := &models.ChangeLogField{
			Field:      "command",
			VType:      reflect.Struct.String(),
			OldVal:     string(old.Command),
			NewVal:     string(new.Command),
			FieldName:  "命令内容",
			OldValName: "",
			NewValName: "",
		}
		if len(g.OldVal.(string)) > 20000 {
			g.OldVal, g.OldValName = "", "数据太大，无法展示"
		}
		if len(g.NewVal.(string)) > 20000 {
			g.OldVal, g.OldValName = "", "数据太大，无法展示"
		}
		content = append(content, g)
	}
	if old.AfterTmpl != new.AfterTmpl {
		content = append(content, &models.ChangeLogField{
			Field:      "after_tmpl",
			VType:      reflect.String.String(),
			OldVal:     old.AfterTmpl,
			NewVal:     new.AfterTmpl,
			FieldName:  "结束模板",
			OldValName: old.AfterTmpl,
			NewValName: new.AfterTmpl,
		})
	}
	if old.Remark != new.Remark {
		content = append(content, &models.ChangeLogField{
			Field:      "remark",
			VType:      reflect.String.String(),
			OldVal:     old.Remark,
			NewVal:     new.Remark,
			FieldName:  "备注",
			OldValName: old.Remark,
			NewValName: new.Remark,
		})
	}
	if old.Status != new.Status {
		content = append(content, &models.ChangeLogField{
			Field:      "status",
			VType:      reflect.Int.String(),
			OldVal:     old.Status,
			NewVal:     new.Status,
			FieldName:  "状态",
			OldValName: models.ConfigStatusMap[old.Status],
			NewValName: models.ConfigStatusMap[new.Status],
		})
	}
	if old.StatusRemark != new.StatusRemark {
		content = append(content, &models.ChangeLogField{
			Field:      "status_remark",
			VType:      reflect.String.String(),
			OldVal:     old.StatusRemark,
			NewVal:     new.StatusRemark,
			FieldName:  "状态描述",
			OldValName: old.StatusRemark,
			NewValName: new.StatusRemark,
		})
	}
	if old.StatusDt != new.StatusDt {
		content = append(content, &models.ChangeLogField{
			Field:      "status_dt",
			VType:      reflect.String.String(),
			OldVal:     old.StatusDt,
			NewVal:     new.StatusDt,
			FieldName:  "状态变更时间",
			OldValName: old.StatusDt,
			NewValName: new.StatusDt,
		})
	}
	if old.MsgSetHash != new.MsgSetHash {
		content = append(content, &models.ChangeLogField{
			Field:      "msg_set",
			VType:      reflect.Struct.String(),
			OldVal:     string(old.MsgSet),
			NewVal:     string(new.MsgSet),
			FieldName:  "消息推送",
			OldValName: "",
			NewValName: "",
		})
	}
	if old.VarFieldsHash != new.VarFieldsHash {
		content = append(content, &models.ChangeLogField{
			Field:      "var_fields",
			VType:      reflect.Struct.String(),
			OldVal:     string(old.VarFields),
			NewVal:     string(new.VarFields),
			FieldName:  "参数变量",
			OldValName: "",
			NewValName: "",
		})
	}
	if old.AfterSleep != new.AfterSleep {
		content = append(content, &models.ChangeLogField{
			Field:      "after_sleep",
			VType:      reflect.Int.String(),
			OldVal:     old.AfterSleep,
			NewVal:     new.AfterSleep,
			FieldName:  "延迟关闭",
			OldValName: strconv.Itoa(old.AfterSleep),
			NewValName: strconv.Itoa(new.AfterSleep),
		})
	}
	if old.ErrRetryNum != new.ErrRetryNum {
		content = append(content, &models.ChangeLogField{
			Field:      "err_retry_num",
			VType:      reflect.Int.String(),
			OldVal:     old.ErrRetryNum,
			NewVal:     new.ErrRetryNum,
			FieldName:  "重试",
			OldValName: strconv.Itoa(old.ErrRetryNum),
			NewValName: strconv.Itoa(new.ErrRetryNum),
		})
	}
	if old.CreateUserId != new.CreateUserId {
		content = append(content, &models.ChangeLogField{
			Field:      "create_user_id",
			VType:      reflect.Int.String(),
			OldVal:     old.CreateUserId,
			NewVal:     new.CreateUserId,
			FieldName:  "创建人",
			OldValName: old.CreateUserName,
			NewValName: new.CreateUserName,
		})
	}
	if old.AuditUserId != new.AuditUserId {
		content = append(content, &models.ChangeLogField{
			Field:      "audit_user_id",
			VType:      reflect.Int.String(),
			OldVal:     old.AuditUserId,
			NewVal:     new.AuditUserId,
			FieldName:  "审核人",
			OldValName: old.AuditUserName,
			NewValName: new.AuditUserName,
		})
	}
	if old.HandleUserIds != new.HandleUserIds {
		content = append(content, &models.ChangeLogField{
			Field:      "handle_user_ids",
			VType:      reflect.Int.String(),
			OldVal:     old.HandleUserIds,
			NewVal:     new.HandleUserIds,
			FieldName:  "处理人",
			OldValName: old.HandleUserNames,
			NewValName: new.HandleUserNames,
		})
	}
	return
}

// 比较并输出差异
func (h *ChangeLogHandle) diffPipeline(old, new *models.CronPipeline) (content []*models.ChangeLogField) {
	content = []*models.ChangeLogField{}
	if old.EntryId != new.EntryId {
		content = append(content, &models.ChangeLogField{
			Field:      "entry_id",
			VType:      reflect.Int.String(),
			OldVal:     old.EntryId,
			NewVal:     new.EntryId,
			FieldName:  "执行编号",
			OldValName: strconv.Itoa(old.EntryId),
			NewValName: strconv.Itoa(new.EntryId),
		})
	}
	if old.Type != new.Type {
		content = append(content, &models.ChangeLogField{
			Field:      "type",
			VType:      reflect.Int.String(),
			OldVal:     old.Type,
			NewVal:     new.Type,
			FieldName:  "类型",
			OldValName: models.ConfigTypeMap[old.Type],
			NewValName: models.ConfigTypeMap[new.Type],
		})
	}
	if old.Name != new.Name {
		content = append(content, &models.ChangeLogField{
			Field:      "name",
			VType:      reflect.String.String(),
			OldVal:     old.Name,
			NewVal:     new.Name,
			FieldName:  "任务名称",
			OldValName: old.Name,
			NewValName: new.Name,
		})
	}
	if old.Spec != new.Spec {
		content = append(content, &models.ChangeLogField{
			Field:      "spec",
			VType:      reflect.String.String(),
			OldVal:     old.Spec,
			NewVal:     new.Spec,
			FieldName:  "执行时间",
			OldValName: old.Spec,
			NewValName: new.Spec,
		})
	}
	if old.VarParams != new.VarParams {
		content = append(content, &models.ChangeLogField{
			Field:      "var_params",
			VType:      reflect.String.String(),
			OldVal:     old.VarParams,
			NewVal:     new.VarParams,
			FieldName:  "变量参数",
			OldValName: old.VarParams,
			NewValName: new.VarParams,
		})
	}
	if string(old.ConfigIds) != string(new.ConfigIds) {
		content = append(content, &models.ChangeLogField{
			Field:      "config_ids",
			VType:      reflect.Slice.String(),
			OldVal:     old.ConfigIds,
			NewVal:     new.ConfigIds,
			FieldName:  "任务",
			OldValName: "", // 任务快照需要简化，先保持空
			NewValName: "",
		})
	}
	if old.ConfigDisableAction != new.ConfigDisableAction {
		content = append(content, &models.ChangeLogField{
			Field:      "config_disable_action",
			VType:      reflect.Int.String(),
			OldVal:     old.ConfigDisableAction,
			NewVal:     new.ConfigDisableAction,
			FieldName:  "任务停用时",
			OldValName: models.DisableActionMap[old.ConfigDisableAction],
			NewValName: models.DisableActionMap[new.ConfigDisableAction],
		})
	}
	if old.ConfigErrAction != new.ConfigErrAction {
		content = append(content, &models.ChangeLogField{
			Field:      "config_err_action",
			VType:      reflect.Int.String(),
			OldVal:     old.ConfigErrAction,
			NewVal:     new.ConfigErrAction,
			FieldName:  "任务错误时",
			OldValName: models.ErrActionMap[old.ConfigErrAction],
			NewValName: models.ErrActionMap[new.ConfigErrAction],
		})
	}
	if old.Interval != new.Interval {
		content = append(content, &models.ChangeLogField{
			Field:      "interval",
			VType:      reflect.Int.String(),
			OldVal:     old.Interval,
			NewVal:     new.Interval,
			FieldName:  "执行间隔",
			OldValName: fmt.Sprintf("%v/秒", old.Interval),
			NewValName: fmt.Sprintf("%v/秒", new.Interval),
		})
	}
	if old.Remark != new.Remark {
		content = append(content, &models.ChangeLogField{
			Field:      "remark",
			VType:      reflect.String.String(),
			OldVal:     old.Remark,
			NewVal:     new.Remark,
			FieldName:  "备注",
			OldValName: old.Remark,
			NewValName: new.Remark,
		})
	}
	if old.Status != new.Status {
		content = append(content, &models.ChangeLogField{
			Field:      "status",
			VType:      reflect.Int.String(),
			OldVal:     old.Status,
			NewVal:     new.Status,
			FieldName:  "状态",
			OldValName: models.ConfigStatusMap[old.Status],
			NewValName: models.ConfigStatusMap[new.Status],
		})
	}
	if old.StatusRemark != new.StatusRemark {
		content = append(content, &models.ChangeLogField{
			Field:      "status_remark",
			VType:      reflect.String.String(),
			OldVal:     old.StatusRemark,
			NewVal:     new.StatusRemark,
			FieldName:  "状态描述",
			OldValName: old.StatusRemark,
			NewValName: new.StatusRemark,
		})
	}
	if old.StatusDt != new.StatusDt {
		content = append(content, &models.ChangeLogField{
			Field:      "status_dt",
			VType:      reflect.String.String(),
			OldVal:     old.StatusDt,
			NewVal:     new.StatusDt,
			FieldName:  "状态变更时间",
			OldValName: old.StatusDt,
			NewValName: new.StatusDt,
		})
	}
	if old.MsgSetHash != new.MsgSetHash {
		content = append(content, &models.ChangeLogField{
			Field:      "msg_set",
			VType:      reflect.Struct.String(),
			OldVal:     string(old.MsgSet),
			NewVal:     string(new.MsgSet),
			FieldName:  "消息推送",
			OldValName: "",
			NewValName: "",
		})
	}
	if old.CreateUserId != new.CreateUserId {
		content = append(content, &models.ChangeLogField{
			Field:      "create_user_id",
			VType:      reflect.Int.String(),
			OldVal:     old.CreateUserId,
			NewVal:     new.CreateUserId,
			FieldName:  "创建人",
			OldValName: old.CreateUserName,
			NewValName: new.CreateUserName,
		})
	}
	if old.AuditUserId != new.AuditUserId {
		content = append(content, &models.ChangeLogField{
			Field:      "audit_user_id",
			VType:      reflect.Int.String(),
			OldVal:     old.AuditUserId,
			NewVal:     new.AuditUserId,
			FieldName:  "审核人",
			OldValName: old.AuditUserName,
			NewValName: new.AuditUserName,
		})
	}
	if old.HandleUserIds != new.HandleUserIds {
		content = append(content, &models.ChangeLogField{
			Field:      "handle_user_ids",
			VType:      reflect.Int.String(),
			OldVal:     old.HandleUserIds,
			NewVal:     new.HandleUserIds,
			FieldName:  "处理人",
			OldValName: old.HandleUserNames,
			NewValName: new.HandleUserNames,
		})
	}
	return
}

// 比较并输出差异
func (h *ChangeLogHandle) diffReceive(old, new *models.CronReceive) (content []*models.ChangeLogField) {
	content = []*models.ChangeLogField{}

	if old.Name != new.Name {
		content = append(content, &models.ChangeLogField{
			Field:      "name",
			VType:      reflect.String.String(),
			OldVal:     old.Name,
			NewVal:     new.Name,
			FieldName:  "任务名称",
			OldValName: old.Name,
			NewValName: new.Name,
		})
	}
	if old.ReceiveTmpl != new.ReceiveTmpl {
		content = append(content, &models.ChangeLogField{
			Field:      "receive_tmpl",
			VType:      reflect.String.String(),
			OldVal:     old.ReceiveTmpl,
			NewVal:     new.ReceiveTmpl,
			FieldName:  "接收模板",
			OldValName: old.ReceiveTmpl,
			NewValName: new.ReceiveTmpl,
		})
	}
	if string(old.ConfigIds) != string(new.ConfigIds) {
		content = append(content, &models.ChangeLogField{
			Field:      "config_ids",
			VType:      reflect.Slice.String(),
			OldVal:     old.ConfigIds,
			NewVal:     new.ConfigIds,
			FieldName:  "任务",
			OldValName: "", // 任务快照需要简化，先保持空
			NewValName: "",
		})
	}
	if old.RuleConfigHash != new.RuleConfigHash {
		g := &models.ChangeLogField{
			Field:      "rule_config",
			VType:      reflect.Struct.String(),
			OldVal:     string(old.RuleConfig),
			NewVal:     string(new.RuleConfig),
			FieldName:  "任务",
			OldValName: "",
			NewValName: "",
		}
		if len(g.OldVal.(string)) > 20000 {
			g.OldVal, g.OldValName = "", "数据太大，无法展示"
		}
		if len(g.NewVal.(string)) > 20000 {
			g.OldVal, g.OldValName = "", "数据太大，无法展示"
		}
		content = append(content, g)
	}

	if old.ConfigDisableAction != new.ConfigDisableAction {
		content = append(content, &models.ChangeLogField{
			Field:      "config_disable_action",
			VType:      reflect.Int.String(),
			OldVal:     old.ConfigDisableAction,
			NewVal:     new.ConfigDisableAction,
			FieldName:  "任务停用时",
			OldValName: models.DisableActionMap[old.ConfigDisableAction],
			NewValName: models.DisableActionMap[new.ConfigDisableAction],
		})
	}
	if old.ConfigErrAction != new.ConfigErrAction {
		content = append(content, &models.ChangeLogField{
			Field:      "config_err_action",
			VType:      reflect.Int.String(),
			OldVal:     old.ConfigErrAction,
			NewVal:     new.ConfigErrAction,
			FieldName:  "任务错误时",
			OldValName: models.ErrActionMap[old.ConfigErrAction],
			NewValName: models.ErrActionMap[new.ConfigErrAction],
		})
	}
	if old.Interval != new.Interval {
		content = append(content, &models.ChangeLogField{
			Field:      "interval",
			VType:      reflect.Int.String(),
			OldVal:     old.Interval,
			NewVal:     new.Interval,
			FieldName:  "执行间隔",
			OldValName: fmt.Sprintf("%v/秒", old.Interval),
			NewValName: fmt.Sprintf("%v/秒", new.Interval),
		})
	}
	if old.Remark != new.Remark {
		content = append(content, &models.ChangeLogField{
			Field:      "remark",
			VType:      reflect.String.String(),
			OldVal:     old.Remark,
			NewVal:     new.Remark,
			FieldName:  "备注",
			OldValName: old.Remark,
			NewValName: new.Remark,
		})
	}
	if old.Status != new.Status {
		content = append(content, &models.ChangeLogField{
			Field:      "status",
			VType:      reflect.Int.String(),
			OldVal:     old.Status,
			NewVal:     new.Status,
			FieldName:  "状态",
			OldValName: models.ConfigStatusMap[old.Status],
			NewValName: models.ConfigStatusMap[new.Status],
		})
	}
	if old.StatusRemark != new.StatusRemark {
		content = append(content, &models.ChangeLogField{
			Field:      "status_remark",
			VType:      reflect.String.String(),
			OldVal:     old.StatusRemark,
			NewVal:     new.StatusRemark,
			FieldName:  "状态描述",
			OldValName: old.StatusRemark,
			NewValName: new.StatusRemark,
		})
	}
	if old.StatusDt != new.StatusDt {
		content = append(content, &models.ChangeLogField{
			Field:      "status_dt",
			VType:      reflect.String.String(),
			OldVal:     old.StatusDt,
			NewVal:     new.StatusDt,
			FieldName:  "状态变更时间",
			OldValName: old.StatusDt,
			NewValName: new.StatusDt,
		})
	}
	if old.MsgSetHash != new.MsgSetHash {
		content = append(content, &models.ChangeLogField{
			Field:      "msg_set",
			VType:      reflect.Struct.String(),
			OldVal:     string(old.MsgSet),
			NewVal:     string(new.MsgSet),
			FieldName:  "消息推送",
			OldValName: "",
			NewValName: "",
		})
	}
	if old.CreateUserId != new.CreateUserId {
		content = append(content, &models.ChangeLogField{
			Field:      "create_user_id",
			VType:      reflect.Int.String(),
			OldVal:     old.CreateUserId,
			NewVal:     new.CreateUserId,
			FieldName:  "创建人",
			OldValName: old.CreateUserName,
			NewValName: new.CreateUserName,
		})
	}
	if old.AuditUserId != new.AuditUserId {
		content = append(content, &models.ChangeLogField{
			Field:      "audit_user_id",
			VType:      reflect.Int.String(),
			OldVal:     old.AuditUserId,
			NewVal:     new.AuditUserId,
			FieldName:  "审核人",
			OldValName: old.AuditUserName,
			NewValName: new.AuditUserName,
		})
	}
	if old.HandleUserIds != new.HandleUserIds {
		content = append(content, &models.ChangeLogField{
			Field:      "handle_user_ids",
			VType:      reflect.Int.String(),
			OldVal:     old.HandleUserIds,
			NewVal:     new.HandleUserIds,
			FieldName:  "处理人",
			OldValName: old.HandleUserNames,
			NewValName: new.HandleUserNames,
		})
	}
	return
}

func (h *ChangeLogHandle) Build() *models.CronChangeLog {
	h.log.CreateDt = time.Now().Format(time.DateTime)
	h.log.CreateUserId = h.user.UserId
	h.log.CreateUserName = h.user.UserName
	// 构建日志
	if h.log.RefType == "config" {
		h.log.RefId = h.conf[1].Id
		content := h.diffConfig(h.conf[0], h.conf[1])
		if len(content) == 0 {
			return nil
		}
		h.log.Content, _ = jsoniter.MarshalToString(content)
		return h.log
	} else if h.log.RefType == "pipeline" {
		h.log.RefId = h.line[1].Id
		content := h.diffPipeline(h.line[0], h.line[1])
		if len(content) == 0 {
			return nil
		}
		h.log.Content, _ = jsoniter.MarshalToString(content)
		return h.log
	} else if h.log.RefType == "receive" {
		h.log.RefId = h.line[1].Id
		content := h.diffReceive(h.rece[0], h.rece[1])
		if len(content) == 0 {
			return nil
		}
		h.log.Content, _ = jsoniter.MarshalToString(content)
		return h.log
	}
	return nil
}

type CronChangeLogData struct {
	db *db.MyDB
}

func NewCronChangeLogData(ctx context.Context) *CronChangeLogData {
	return &CronChangeLogData{
		db: db.New(ctx),
	}
}

// 写入数据
func (m *CronChangeLogData) Write(h *ChangeLogHandle) error {
	g := h.Build()
	if g == nil {
		return nil
	}
	return m.db.Create(g).Error
}

func (m *CronChangeLogData) ListPage(where *db.Where, page, size int, list interface{}) (total int64, err error) {
	str, args := where.Build()

	return m.db.Paginate(list, page, size, "cron_change_log", "*", "id desc", str, args...)
}
