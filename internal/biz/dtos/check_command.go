package dtos

import (
	"cron/internal/basic/enum"
	"cron/internal/basic/grpcurl"
	"cron/internal/models"
	"cron/internal/pb"
	"errors"
	"fmt"
	"regexp"
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
	if http.Timeout < 0 {
		http.Timeout = 0
	}
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
	if sql.Origin == enum.SqlStatementSourceGit && sql.GitSourceId == 0 && len(sql.Statement) > 0 {
		return fmt.Errorf("请选择 git 连接")
	}
	for _, item := range sql.Statement {
		if sql.Origin != item.Type {
			continue // 来源不一致的忽略
		}
		if sql.Origin == enum.SqlStatementSourceLocal {
			if item.Local == "" {
				return errors.New("未设置 sql 执行语句")
			}
		} else if sql.Origin == enum.SqlStatementSourceGit {
			//if item.Git.LinkId == 0 {
			//	return errors.New("未设置 sql 语句 连接")
			//}
			if item.Git.Owner == "" {
				return errors.New("未设置 sql 语句 仓库空间")
			}
			if !regexp.MustCompile(`^[a-zA-Z][\w-]{1,}[a-zA-Z0-9]$`).MatchString(item.Git.Owner) {
				return errors.New("仓库空间: 只允许字母、数字或者下划线（_）、中划线（-），至少 2 个字符，必须以字母开头，不能以特殊字符结尾")
			}
			if item.Git.Project == "" {
				return errors.New("未设置 sql 语句 项目名称")
			}
			if len(item.Git.Path) < 1 {
				return errors.New("未设置 sql 语句 文件路径")
			}
		} else {
			return errors.New("sql来源有误")
		}
		if _, ok := enum.BoolMap[item.IsBatch]; !ok {
			return errors.New("请确认是否批量解析")
		}
	}
	if _, ok := enum.SqlDriverMap[sql.Driver]; !ok {
		return errors.New("sql 驱动设置有误")
	}

	if _, ok := models.SqlErrActionMap[sql.ErrAction]; !ok {
		return errors.New("未设置 sql 错误行为")
	}
	if sql.ErrAction == models.SqlErrActionRollback && sql.Interval > 0 {
		return errors.New("事务回滚 时禁用 执行间隔")
	}
	if sql.Interval < 0 {
		return errors.New("sql 执行间隔不得小于0")
	}
	return nil
}

func CheckCmd(cmd *pb.CronCmd) error {
	if cmd.Type == "" && cmd.Host.Id <= 0 { // 远程主机可以不用指定type
		return fmt.Errorf("未指定命令行类型")
	}
	if cmd.Host.Id != -1 && cmd.Host.Id <= 0 { // -1.本机
		return fmt.Errorf("主机选择有误")
	}
	if cmd.Origin == enum.SqlStatementSourceLocal {
		if cmd.Statement.Local == "" {
			return fmt.Errorf("请输入 cmd 命令类容")
		}
	} else if cmd.Origin == enum.SqlStatementSourceGit {
		if cmd.Statement.Git.LinkId == 0 {
			return errors.New("未设置 命令 连接")
		}
		if cmd.Statement.Git.Owner == "" {
			return errors.New("未设置 命令 仓库空间")
		}
		if cmd.Statement.Git.Project == "" {
			return errors.New("未设置 命令 项目名称")
		}
		pathLen := len(cmd.Statement.Git.Path)
		if pathLen == 0 {
			return errors.New("未设置 命令 文件路径")
		} else if pathLen > 1 {
			return errors.New("命令 文件路径 不支持多文件")
		}
	} else {
		return fmt.Errorf("未指定命令行来源")
	}

	return nil
}

func CheckJenkins(jks *pb.CronJenkins) error {
	if jks.Source == nil || jks.Source.Id == 0 {
		return fmt.Errorf("未选择链接")
	}
	if jks.Name == "" {
		return fmt.Errorf("项目名称不得为空")
	}
	pl := len(jks.Params)
	if jks.ParamsMode == models.ParamModeDefault {
		for i, param := range jks.Params {
			if param.Key == "" && i < (pl-1) {
				return fmt.Errorf("参数 %v 名称不得为空", i+1)
			}
		}
	} else if jks.ParamsMode == models.ParamModeGroup {
		if len(jks.ParamsGroup) == 0 {
			return fmt.Errorf("至少添加一个参数组")
		}
		for i, group := range jks.ParamsGroup {
			pl := len(group.Params)
			for j, param := range group.Params {
				if param.Key == "" && j < (pl-1) {
					return fmt.Errorf("参数组 %v 第 %v 个参数 名称不得为空", i+1, j+1)
				}
			}
		}
	} else {
		return fmt.Errorf("参数模式有误")
	}
	return nil
}

func CheckGit(raw, c *pb.CronGit) error {
	if c.LinkId <= 0 {
		return fmt.Errorf("未指定有效连接")
	}
	if len(c.Events) == 0 {
		return fmt.Errorf("未指定事件")
	}
	for i, e := range c.Events {
		switch e.Id {
		case enum.GitEventPullsDetail:
			if e.PRDetail.Owner == "" {
				return errors.New("git 仓库空间 未设置")
			}
			if !regexp.MustCompile(`^[a-zA-Z][\w-]{1,}[a-zA-Z0-9]$`).MatchString(e.PRDetail.Owner) {
				return errors.New("git 仓库空间 只允许字母、数字或者下划线（_）、中划线（-），至少 2 个字符，必须以字母开头，不能以特殊字符结尾")
			}
			if e.PRDetail.Repo == "" {
				return errors.New("git 项目名称 未设置")
			}
			if e.PRDetail.Number == "" {
				return errors.New("git 仓库PR的序数为必填")
			}
		case enum.GitEventPullsIsMerge:
			if e.PRIsMerge.Owner == "" {
				return errors.New("git 仓库空间 未设置")
			}
			if !regexp.MustCompile(`^[a-zA-Z][\w-]{1,}[a-zA-Z0-9]$`).MatchString(e.PRIsMerge.Owner) {
				return errors.New("git 仓库空间 只允许字母、数字或者下划线（_）、中划线（-），至少 2 个字符，必须以字母开头，不能以特殊字符结尾")
			}
			if e.PRIsMerge.Repo == "" {
				return errors.New("git 项目名称 未设置")
			}
			if e.PRIsMerge.Number == "" {
				return errors.New("git 仓库PR的序数为必填")
			}
			if e.PRIsMerge.State == "" {
				return errors.New("git 目标状态 必填")
			}
			if e.PRIsMerge.State != "open" && e.PRIsMerge.State != "merge" {
				return errors.New("git 目标状态 错误")
			}
		case enum.GitEventPullsMerge:
			if e.PRMerge.Owner == "" {
				return errors.New("git 仓库空间 未设置")
			}
			if !regexp.MustCompile(`^[a-zA-Z][\w-]{1,}[a-zA-Z0-9]$`).MatchString(e.PRMerge.Owner) {
				return errors.New("git 仓库空间 只允许字母、数字或者下划线（_）、中划线（-），至少 2 个字符，必须以字母开头，不能以特殊字符结尾")
			}
			if e.PRMerge.Repo == "" {
				return errors.New("git 项目名称 未设置")
			}
			if e.PRMerge.Number == "" {
				return errors.New("git 仓库PR的序数为必填")
			}
			if e.PRMerge.MergeMethod == "" {
				return errors.New("git 合并方式不得为空")
			}
		case enum.GitEventFileUpdate:
			if e.FileUpdate.Owner == "" {
				return errors.New("git 仓库空间 未设置")
			}
			if !regexp.MustCompile(`^[a-zA-Z][\w-]{1,}[a-zA-Z0-9]$`).MatchString(e.FileUpdate.Owner) {
				return errors.New("git 仓库空间 只允许字母、数字或者下划线（_）、中划线（-），至少 2 个字符，必须以字母开头，不能以特殊字符结尾")
			}
			if e.FileUpdate.Repo == "" {
				return errors.New("git 项目名称 未设置")
			}
			if e.FileUpdate.Path == "" {
				return errors.New("git 文件路径为必填")
			}
			if e.FileUpdate.Content == "" {
				return errors.New("git 文件内容不得为空")
			}
			if e.FileUpdate.Message == "" {
				return errors.New("git 提交描述不得为空")
			}
		default:
			return fmt.Errorf("未支持的事件 %v-%v", i, e.Id)
		}
	}
	return nil
}
