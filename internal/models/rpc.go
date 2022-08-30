package models

type GrpcRequest struct {
	body string
}

// grpc必须实现
func (m *GrpcRequest) Reset() {
	m.body = ""
}

// grpc必须实现
func (m *GrpcRequest) String() string {
	return m.body
}

// grpc必须实现
func (m *GrpcRequest) ProtoMessage() {

}

func (m *GrpcRequest) SetParam(param string) {
	m.body = param
}
