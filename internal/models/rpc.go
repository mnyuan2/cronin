package models

type GrpcRequest struct {
	Body string `protobuf:"bytes,3,opt,name=body"`
}

// grpc必须实现
func (m *GrpcRequest) Reset() {
	m.Body = ""
}

// grpc必须实现
func (m *GrpcRequest) String() string {
	return m.Body
}

// grpc必须实现
func (m *GrpcRequest) ProtoMessage() {

}

func (m *GrpcRequest) SetParam(param string) {
	m.Body = param
}
