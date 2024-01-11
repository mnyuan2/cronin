package grpcurl

import (
	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/desc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// 请求处理
type MyEventHandler struct {
	formatter    Formatter
	Method       string
	ReqHeaders   map[string][]string
	ReqMessages  []byte
	RespHeaders  map[string][]string
	RespMessages string
	RespTrailers map[string][]string // ?
	status       *status.Status
}

func NewMyEventHandler(formatter Formatter) *MyEventHandler {
	return &MyEventHandler{formatter: formatter}
}

// OnResolveMethod与正在被调用的方法的描述符一起调用。
func (h *MyEventHandler) OnResolveMethod(md *desc.MethodDescriptor) {
	h.Method, _ = GetDescriptorText(md, nil)
}

// 使用正在发送的请求元数据调用OnSendHeaders。
func (h *MyEventHandler) OnSendHeaders(md metadata.MD) {
	h.ReqHeaders = md
}

// OnReceiveHeaders在接收到响应头时被调用。
func (h *MyEventHandler) OnReceiveHeaders(md metadata.MD) {
	h.RespHeaders = md
}

// 每收到一个响应消息就调用
func (h *MyEventHandler) OnReceiveResponse(resp proto.Message) {
	//proto.Size(resp)
	h.RespMessages, _ = h.formatter(resp)
	// 格式化响应消息失败
}

// OnReceiveTrailers在接收到响应拖车和最终RPC状态时调用。
func (h *MyEventHandler) OnReceiveTrailers(stat *status.Status, md metadata.MD) {
	h.status = stat
	h.RespTrailers = md
}

func (h *MyEventHandler) SetStatus(stat *status.Status) {
	h.status = stat
}

func (h *MyEventHandler) GetStatus() *status.Status {
	return h.status
}
