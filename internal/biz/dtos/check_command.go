package dtos

import (
	"cron/internal/basic/grpcurl"
	"cron/internal/models"
	"cron/internal/pb"
	"errors"
	"fmt"
	"strings"
)

// http 设置校验
func CheckHttp(http *pb.CronHttp) error {
	if !strings.HasPrefix(http.Url, "http://") && !strings.HasPrefix(http.Url, "https://") {
		return fmt.Errorf("请输入 http:// 或 https:// 开头的规范地址")
	}
	if http.Method == "" {
		return errors.New("请输入请求method")
	}
	if models.ProtocolHttpMethodMap()[http.Method] == "" {
		return errors.New("未支持的请求method")
	}
	//if http.Body != "" {
	//	temp := map[string]any{}
	//	if err := jsoniter.UnmarshalFromString(http.Body, &temp); err != nil {
	//		return fmt.Errorf("http body 输入不规范，请确认json字符串是否规范")
	//	}
	//}
	return nil
}

func CheckRPC(rpc *pb.CronRpc) error {
	if rpc.Method != "GRPC" {
		return fmt.Errorf("rpc 请选择请求模式")
	}
	if rpc.Proto == "" {
		return fmt.Errorf("rpc 请完善proto文件内容")
	}
	if rpc.Addr == "" {
		return fmt.Errorf("rpc 请完善请求地址")
	}
	if rpc.Action == "" {
		return fmt.Errorf("rpc 请完善请求方法")
	}
	fds, err := grpcurl.ParseProtoString(rpc.Proto)
	if err != nil {
		return err
	}
	rpc.Actions = grpcurl.ParseProtoMethods(fds)
	actionOk := false
	for _, item := range rpc.Actions {
		if item == rpc.Action {
			actionOk = true
		}
	}
	if !actionOk {
		return fmt.Errorf("rpc 请求方法与proto不符")
	}
	return nil
}

func CheckSql(sql *pb.CronSql) error {
	if sql.Source.Id == 0 {
		return fmt.Errorf("请选择 sql 连接")
	}
	if len(sql.Statement) == 0 {
		return errors.New("未设置 sql 执行语句")
	}
	name, ok := models.SqlErrActionMap[sql.ErrAction]
	if !ok {
		return errors.New("未设置 sql 错误行为")
	}
	sql.ErrActionName = name
	if sql.ErrAction == models.SqlErrActionRollback && sql.Interval > 0 {
		return errors.New("事务回滚 时禁用 执行间隔")
	}
	if sql.Interval < 0 {
		sql.Interval = 0
	}
	return nil
}
