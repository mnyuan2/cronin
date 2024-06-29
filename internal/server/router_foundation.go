package server

import (
	"cron/internal/biz"
	"cron/internal/pb"
	"github.com/gin-gonic/gin"
	"gopkg.in/antage/eventsource.v1"
)

var es = eventsource.New(nil, nil)

//defer es.Close()

// 获取选项列表
func routerDicGets(ctx *gin.Context) {
	r := &pb.DicGetsRequest{}
	if err := ctx.BindQuery(r); err != nil {
		NewReply(ctx).SetError(pb.ParamError, err.Error()).RenderJson()
		return
	}
	user, err := GetUser(ctx)
	if err != nil {
		NewReply(ctx).SetError(pb.UserNotExist, err.Error()).RenderJson()
		return
	}
	rep, err := biz.NewDicService(ctx.Request.Context(), user).DicGets(r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

func routerSystemInfo(ctx *gin.Context) {

	rep, err := biz.NewDicService(ctx.Request.Context(), nil).SystemInfo(nil)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

// 解析proto文件
func routerParseProto(ctx *gin.Context) {
	r := &pb.ParseProtoRequest{}
	if err := ctx.BindJSON(r); err != nil {
		NewReply(ctx).SetError(pb.ParamError, err.Error()).RenderJson()
		return
	}
	rep, err := biz.NewDicService(ctx.Request.Context(), nil).ParseProto(r)
	NewReply(ctx).SetReply(rep, err).RenderJson()
}

func routerEvents(ctx *gin.Context) {
	// 设置HTTP头，声明内容类型为text/event-stream
	//ctx.Header("Content-Type", "text/event-stream")
	//ctx.Header("Cache-Control", "no-cache")
	// 使用循环来向客户端发送数据
	//for {
	//	// 构造sse数据格式，包括事件类型和数据内容
	//	eventData := fmt.Sprintf("data: %s\n\n", time.Now().Format(time.RFC3339))
	//
	//	// 将数据写入ResponseWriter
	//	_, err := ctx.Writer.Write([]byte(eventData))
	//	if err != nil {
	//		fmt.Println("Error writing SSE data: ", err)
	//		break
	//	}
	//
	//	// 强制将数据刷新到客户端，保证客户端可以及时接收数据
	//	ctx.Writer.Flush()
	//
	//	// 模拟每秒发送一次数据
	//	time.Sleep(1 * time.Second)
	//}
	// 设置响应头，声明内容类型为 text/event-stream
	es.ServeHTTP(ctx.Writer, ctx.Request)
}
