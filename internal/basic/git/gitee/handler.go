package gitee

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"net/http"
	"time"
)

// 处理结果
type Handler struct {
	ctx            context.Context
	startTime      time.Time
	endTime        time.Time
	General        *General
	RequestHeader  http.Header
	RequestBody    []byte
	ResponseHeader http.Header
	ResponseBody   []byte
}

type General struct {
	Method     string
	Url        string
	StatusCode int
}

func NewHandler(ctx context.Context) *Handler {
	return &Handler{
		ctx:     ctx,
		General: &General{},
	}
}

// 通用信息
func (m *Handler) OnGeneral(method, url string, statusCode int) {
	m.General.Method = method
	m.General.Url = url
	m.General.StatusCode = statusCode
}

// 请求头
func (m *Handler) OnRequestHeader(header http.Header) {
	m.RequestHeader = header
}

// 请求body
func (m *Handler) OnRequestBody(b []byte) {
	m.RequestBody = b
}

// 响应头
func (m *Handler) OnResponseHeader(header http.Header) {
	m.ResponseHeader = header
}

// 响应内容
func (m *Handler) OnResponseBody(b []byte) {
	m.ResponseBody = b
}

// 上下文
func (m *Handler) GetContext() context.Context {
	return m.ctx
}

func (m *Handler) RequestHeaderBytes() []byte {
	if m.RequestHeader == nil {
		return nil
	}
	b, _ := jsoniter.Marshal(m.RequestHeader)
	return b
}

func (m *Handler) ResponseHeaderBytes() []byte {
	if m.ResponseHeader == nil {
		return nil
	}
	b, _ := jsoniter.Marshal(m.ResponseHeader)
	return b
}

func (m *Handler) StartTime() time.Time {
	return m.startTime
}

func (m *Handler) EndTime() time.Time {
	return m.endTime
}

// 返回字符串表示
func (m *Handler) String() string {
	s, _ := jsoniter.MarshalToString(map[string]any{
		"general":        m.General,
		"request_header": m.RequestHeader,
		"RequestBody":    string(m.RequestBody),
		"ResponseHeader": m.ResponseHeader,
		"ResponseBody":   string(m.ResponseBody),
	})
	return s
}
