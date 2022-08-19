package server

import (
	"cron/internal/pb"
	"fmt"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"
)

// 响应结构体
type GinReply struct {
	Code    string       `json:"code"`
	Message string       `json:"message"`
	Data    interface{}  `json:"data"`
	b       []byte       `json:"-"`
	c       *gin.Context `json:"-"`
}

func NewReply(ctx *gin.Context) *GinReply {
	return &GinReply{c: ctx}
}

// 响应数据写入
func (r *GinReply) Render(w http.ResponseWriter) (err error) {
	if r.b, err = jsoniter.Marshal(r); err != nil {
		log.Println("响应内容转义异常", err)
		r.Code = pb.SysError
		r.b = []byte(`{"code":2001, "message":"响应内容转义异常:` + err.Error() + `", "data":null}`)
	}

	_, err = w.Write(r.b)
	return err
}

// 统一响应头
func (r *GinReply) WriteContentType(w http.ResponseWriter) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = []string{"application/json; charset=utf-8"}
	} else if val[0] == "text/plain; charset=utf-8" {
		header["Content-Type"][0] = "application/json; charset=utf-8" // 响应的普通文本头，替换为json头。
	}
}

func (r *GinReply) SetError(Code string, depicts ...string) *GinReply {
	msg, ok := pb.CodeList[Code]
	if !ok {
		panic("错误码未定义！")
	}
	if len(depicts) > 0 {
		msg += fmt.Sprintf("[ %v ]", strings.Join(depicts, ","))
	}

	r.Code = Code
	r.Message = msg
	return r
}

func (r *GinReply) SetSuccess(data ...interface{}) *GinReply {
	r.Code = pb.Success
	r.Message = pb.CodeList[r.Code]
	if len(data) > 0 {
		r.Data = data[0]
	}

	return r
}

// SetReply 设置响应结果
//
// 快捷方法， 内部集成了成功和错误的响应流程。
func (r *GinReply)SetReply(resp interface{}, err error)*GinReply{
	if err != nil {
		r.SetError(pb.OperationFailure, err.Error())
	}else {
		r.SetSuccess(resp)
	}
	return r
}

// 渲染结果
func (r *GinReply) RenderJson() {
	r.c.Render(http.StatusOK, r)
	if r.Code == pb.Success {
		r.c.Set("result_status", "success")
	} else {
		r.c.Set("result_status", "fail")
	}
	r.c.Set("result_body", r.b)
}

// 响应文件
func (r *GinReply) RenderFile(fileName string, data []byte) {
	if len(data) <= 0 {
		r.SetError(pb.NotFound, "没有可下载的数据").RenderJson()
		return
	}

	r.c.Header("Content-Type", "application/octet-stream")
	r.c.Header("Content-Disposition", fmt.Sprintf("attachment; filename*=utf-8''%s", url.QueryEscape(path.Base(fileName))))
	if _, err := r.c.Writer.Write(data); err != nil {
		r.SetError(pb.OperationFailure, "文件下载错误", err.Error()).RenderJson()
	} else {
		r.c.Set("response", fileName)
	}
}

// 解析data到指定结构体
func (r *GinReply) UnmarshalData(data interface{}) error {
	b, e := jsoniter.Marshal(r.Data)
	if e != nil {
		return e
	}
	return jsoniter.Unmarshal(b, data)
}

// 渲染结果
func Render(c *gin.Context, re *GinReply) {
	c.Render(http.StatusOK, re)
	if re.Code == pb.Success {
		c.Set("result_status", "success")
	} else {
		c.Set("result_status", "fail")
	}
	c.Set("result_body", re.b)
	return
}
