package sse

import (
	"gopkg.in/antage/eventsource.v1"
	"log"
	"net/http"
	"sync"
	"time"
)

// 服务启动后，全局唯一
var se *sse
var once sync.Once

func init() {
	// 创建一个新的 EventSource 实例
	once.Do(func() {
		se = &sse{
			es: eventsource.New(nil, nil),
		}
	})
}

type sse struct {
	es eventsource.EventSource
}

// 获取服务
func Serve() *sse {
	return se
}

func (m *sse) Close() {
	m.es.Close()
}

// 在后台每秒向所有连接的客户端发送当前时间
func (m *sse) ListenHandler(event string, interval time.Duration, handler func() string) {
	for {
		data := handler()
		m.es.SendEventMessage(data, event, "") // 发送消息
		log.Printf("在线人数: %d", m.es.ConsumersCount())
		time.Sleep(interval) // 间隔时间
	}
}

// 发送消息
func (m *sse) SendEventMessage(event, data string) {
	m.es.SendEventMessage(data, event, "")
}

// SSE 路由，处理客户端的 SSE 连接
func (m *sse) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	//ctx.Header("Content-Type", "text/event-stream")
	//ctx.Header("Cache-Control", "no-cache")
	//ctx.Header("Connection", "keep-alive")
	//ctx.Header("Access-Control-Allow-Origin", "*")
	m.es.ServeHTTP(writer, request)
}
