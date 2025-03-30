package models

import (
	"bytes"
	"cron/internal/basic/config"
	"cron/internal/basic/db"
	"cron/internal/basic/enum"
	"cron/internal/pb"
	"crypto/md5"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"log"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// mysql最低版本
// 5.7.8 开始支持json函数，低于程序会报错
var mysqlLower = []int{5, 7, 7}

// 注册表结构
func AutoMigrate(Db *db.MyDB) {
	if config.DbConf().Driver == db.DriverMysql {
		if err := mysqlLowerCheck(Db); err != nil {
			panic(err.Error())
		}
		// 迁移表结构
		err := Db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4").
			AutoMigrate(&CronSetting{}, &CronConfig{}, &CronPipeline{}, &CronReceive{}, &CronLogSpan{}, &CronLogSpanIndex{},
				&CronUser{}, &CronAuthRole{}, &CronChangeLog{}, &CronTag{})
		if err != nil {
			panic(fmt.Sprintf("mysql 表初始化失败，%s", err.Error()))
		}
	} else if config.DbConf().Driver == db.DriverSqlite {
		err := Db.AutoMigrate(&CronSetting{}, &CronConfig{}, &CronPipeline{}, &CronReceive{}, &CronLogSpan{}, &CronLogSpanIndex{},
			&CronUser{}, &CronAuthRole{}, &CronChangeLog{}, &CronTag{})
		if err != nil {
			panic(fmt.Sprintf("mysql 表初始化失败，%s", err.Error()))
		}
	} else {
		panic(fmt.Sprintf("驱动异常，%s", config.DbConf().Driver))
	}

	ti := time.Now()
	// 初始化数据
	err := Db.Where("scene=? and status=?", SceneEnv, enum.StatusActive).FirstOrCreate(&CronSetting{
		Scene:    "env",
		Name:     "public",
		Title:    "public",
		Content:  `{"default":2}`,
		Status:   enum.StatusActive,
		CreateDt: ti.Format(time.DateTime),
		UpdateDt: ti.Format(time.DateTime),
	}).Error
	if err != nil {
		panic(fmt.Sprintf("cron_setting 表默认行数据初始化失败，%s", err.Error()))
	}
	msg := &CronSetting{}
	err = Db.Where("scene=?", SceneMsg).Find(msg).Error
	if err != nil {
		panic(fmt.Sprintf("cron_setting 表默认行数据初始化失败，%s", err.Error()))
	}
	if msg.Id == 0 { // 后期会有多条默认消息模板
		Db.CreateInBatches([]*CronSetting{
			{
				Scene: "msg",
				Name:  "",
				Title: "企微xx群",
				Content: `{
    "http": {
        "method": "POST",
        "url": "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=xx",
        "body": "{\n    \"msgtype\": \"text\",\n    \"text\": {\n        \"content\": \"时间：[[.log.create_dt]]\\n任务：【[[.env]]】[[.config.name]]\\n状态：[[.log.status_name]]-[[.log.status_desc]][[if gt .log.retry_number 0]] 【重试 [[.log.retry_number]]】[[end]]\\n耗时：[[.log.duration]]秒\\n响应：[[.log.body]]\",\n        \"mentioned_mobile_list\": [[.user.mobile]]\n    }\n}",
        "header": []
    }
}`,
				Status:   enum.StatusActive,
				CreateDt: ti.Format(time.DateTime),
				UpdateDt: ti.Format(time.DateTime),
			},
		}, 10)
	}

	// 初始化角色
	role := &CronAuthRole{}
	err = Db.Find(role).Error
	if err != nil {
		panic(fmt.Sprintf("cron_auth_role 表默认行数据初始化失败，%s", err.Error()))
	}
	if role.Id == 0 {
		Db.CreateInBatches([]*CronAuthRole{
			{
				Id:      1,
				Name:    "管理员",
				Remark:  "所有权限",
				AuthIds: "20,21,22,23,24,25,30,31,32,33,34,35,60,61,62,63,70,71,72,74,75,80,81,82,83,90,91,92,95,100,101,102,104,105,120,121,132,133,150,151,152,153,154,155,160,161,162,163",
				Status:  enum.StatusActive,
			},
			{
				Id:      2,
				Name:    "负责人",
				Remark:  "负责任务的创建与审核",
				AuthIds: "20,21,22,23,24,25,30,31,32,33,34,35,61,71,80,81,82,83,91,151,152,153,155,160,161,162,163",
				Status:  enum.StatusActive,
			},
			{
				Id:      3,
				Name:    "标准",
				Remark:  "负责任务的创建并提交审核",
				AuthIds: "21,22,23,25,31,32,33,35,71,81,151,155,161",
				Status:  enum.StatusActive,
			},
		}, 10)
	}

	// 初始化超管用户
	user := &CronUser{}
	err = Db.Where("role_ids !=''").Find(user).Error
	if err != nil {
		panic(fmt.Sprintf("cron_user 表默认行数据初始化失败，%s", err.Error()))
	}
	if user.Id == 0 {
		password, err := SqlSourceEncrypt("123456")
		if err != nil {
			panic("请确认秘钥是否正确设置，" + err.Error())
		}
		err = Db.CreateInBatches([]*CronUser{{
			Account:  "ROOT",
			Username: "超管",
			Status:   enum.StatusActive,
			UpdateDt: ti.Format(time.DateTime),
			CreateDt: ti.Format(time.DateTime),
			Password: password,
			RoleIds:  "1",
		}}, 10).Error
		if err != nil {
			panic("超管账户初始化异常，" + err.Error())
		}
	}

	// 历史 数据修正
	historyDataRevise(Db)
}

// mysql 最低版本检测
func mysqlLowerCheck(db *db.MyDB) error {
	version := ""
	err := db.Raw("SELECT VERSION()").Scan(&version).Error
	if err != nil {
		return fmt.Errorf("mysql 版本获取失败，%s", err.Error())
	}

	temp1 := strings.Split(version, "-")
	temp2 := strings.Split(temp1[0], ".")
	isLower := true
	for i, n := range temp2 {
		val, _ := strconv.Atoi(n)
		if mysqlLower[i] < val {
			isLower = false
			break
		}
	}
	if isLower {
		return fmt.Errorf("mysql最低要求版本 5.7.8 当前为 %s", version)
	}
	return nil
}

// 历史源修正
// 解决 0.6.1 之前的版本格式不一致问题
func historyDataRevise(db *db.MyDB) {
	set := &CronSetting{}
	db.Where("scene='sys_tag_history_update'").Find(set)
	if set.Scene == "" {
		// sql_source 历史数据修正
		if err := db.Exec(`UPDATE cron_setting set content=concat('{"sql":',content,'}') WHERE scene='sql_source' and content->'$.sql' is null;`).Error; err != nil {
			log.Println("历史 sql_source 数据修正错误", err.Error())
		}

		// config cmd 历史数据修正
		cmdType := "sh"
		if runtime.GOOS == "windows" {
			cmdType = "cmd"
		}
		err := db.Exec(fmt.Sprintf(`UPDATE cron_config SET command=JSON_REPLACE(command,'$.cmd', CAST(concat('{"type":"%s","origin":"local","statement":{"type":"local","git":{},"local":',command->'$.cmd','}}') as JSON)) WHERE JSON_TYPE(command->'$.cmd') = 'STRING';`, cmdType)).Error
		if err != nil {
			log.Println("历史 config cmd 数据修正错误", err.Error())
		}
		// config sql 历史数据修正
		list := []*CronConfig{}
		db.Where("JSON_TYPE(command->'$.sql') = 'OBJECT' and command->'$.sql.origin' is null").Select("id", "command").Find(&list)
		if len(list) > 0 {
			type CronSql struct {
				Statement []string `json:"statement"` // sql语句多条
			}
			type CronConfigCommand struct {
				Sql *CronSql `json:"sql"`
			}
			for _, item := range list {
				cmd := &CronConfigCommand{Sql: &CronSql{Statement: []string{}}}
				if er := jsoniter.Unmarshal(item.Command, cmd); err != nil {
					log.Println("	sql 解析错误", item.Id, er.Error())
					continue
				}
				newStatement := make([]map[string]string, len(cmd.Sql.Statement))
				for i, statement := range cmd.Sql.Statement {
					newStatement[i] = map[string]string{
						"type":  "local",
						"local": statement,
					}
				}
				str, _ := jsoniter.MarshalToString(newStatement)
				updateSql := `UPDATE cron_config set command=JSON_SET(command, '$.sql.origin', 'local', '$.sql.statement', CAST(? as JSON)) WHERE id=?`
				if er := db.Exec(updateSql, str, item.Id).Error; er != nil {
					log.Println("	sql 修正错误: ", updateSql, er.Error())
				}
			}
		}
		set.Scene = "sys_tag_history_update"
		set.Content = `{"version":"0.6.1"}`
		db.Create(set)
	}

	//jso
	if set.Content == `{"version":"0.6.1"}` {
		// 新增了sql驱动字段，对没有值的历史数据进行初始化。
		list := []*CronSetting{}
		db.Where("scene='sql_source'").Find(&list)
		for _, item := range list {
			source := &pb.SettingSource{
				Sql:     &pb.SettingSqlSource{},
				Git:     &pb.SettingGitSource{},
				Jenkins: &pb.SettingJenkinsSource{},
				Host:    &pb.SettingHostSource{},
			}
			if er := jsoniter.UnmarshalFromString(item.Content, source); er == nil {
				if source.Sql.Driver == "" {
					source.Sql.Driver = enum.SqlDriverMysql
					item.Content, er = jsoniter.MarshalToString(source)
					if er == nil {
						db.Select("content").Updates(item)
					}
				}
			}
		}
		set.Content = `{"version":"0.7.0"}`
		db.Select("content").Updates(set)
	}
	if set.Content == `{"version":"0.7.0"}` { // 升级所有消息状态从数值升级到数组；消息模板升级为 text/template 解析
		type CronMsgSet struct {
			MsgId         int   `json:"msg_id"`
			Status        int   `json:"status"`
			NotifyUserIds []int `json:"notify_user_ids"`
		}
		// 接收任务
		receiveList := []*CronReceive{}
		db.Find(&receiveList)
		for _, item := range receiveList {
			if len(item.MsgSet) < 5 {
				continue
			}
			msgs := []*CronMsgSet{}
			newMsgs := []*pb.CronMsgSet{}
			if err := jsoniter.Unmarshal(item.MsgSet, &msgs); err != nil {
				continue
			}
			for _, m := range msgs {
				newMsgs = append(newMsgs, &pb.CronMsgSet{MsgId: m.MsgId, Status: []int{m.Status}, NotifyUserIds: m.NotifyUserIds})
			}
			item.MsgSet, _ = jsoniter.Marshal(newMsgs)
			item.MsgSetHash = fmt.Sprintf("%x", md5.Sum(item.MsgSet))
			db.Select("msg_set", "msg_set_hash").Updates(item)
		}
		// 流水线任务
		pipelineList := []*CronPipeline{}
		db.Find(&pipelineList)
		for _, item := range pipelineList {
			if len(item.MsgSet) < 5 {
				continue
			}
			msgs := []*CronMsgSet{}
			newMsgs := []*pb.CronMsgSet{}
			if err := jsoniter.Unmarshal(item.MsgSet, &msgs); err != nil {
				continue
			}
			for _, m := range msgs {
				newMsgs = append(newMsgs, &pb.CronMsgSet{MsgId: m.MsgId, Status: []int{m.Status}, NotifyUserIds: m.NotifyUserIds})
			}
			item.MsgSet, _ = jsoniter.Marshal(newMsgs)
			item.MsgSetHash = fmt.Sprintf("%x", md5.Sum(item.MsgSet))
			db.Select("msg_set", "msg_set_hash").Updates(item)
		}
		// 任务
		configList := []*CronConfig{}
		db.Find(&configList)
		for _, item := range configList {
			if len(item.MsgSet) < 5 {
				continue
			}
			msgs := []*CronMsgSet{}
			newMsgs := []*pb.CronMsgSet{}
			if err := jsoniter.Unmarshal(item.MsgSet, &msgs); err != nil {
				continue
			}
			for _, m := range msgs {
				newMsgs = append(newMsgs, &pb.CronMsgSet{MsgId: m.MsgId, Status: []int{m.Status}, NotifyUserIds: m.NotifyUserIds})
			}
			item.MsgSet, _ = jsoniter.Marshal(newMsgs)
			item.MsgSetHash = fmt.Sprintf("%x", md5.Sum(item.MsgSet))
			db.Select("msg_set", "msg_set_hash").Updates(item)
		}

		// 消息模板更新
		msgsArgs := []string{
			"env",
			"config.name",
			"config.protocol_name",
			"log.status_name",
			"log.status_desc",
			"log.body",
			"log.duration",
			"log.create_dt",
			"user.username",
			"user.mobile",
		}
		msgs := []*CronSetting{}
		db.Where("scene=?", SceneMsg).Find(&msgs)
		for _, msg := range msgs {
			b := []byte(msg.Content)
			for _, k := range msgsArgs {
				b = bytes.Replace(b, []byte("[["+k+"]]"), []byte("[[."+k+"]]"), -1)
			}
			msg.Content = string(b)
			db.Select("content").Updates(msg)
		}

		set.Content = `{"version":"0.8.1"}`
		db.Select("content").Updates(set)
	}
	if set.Content == `{"version":"0.8.1"}` {
		list := []*CronSetting{}
		db.Where("scene in ?", []string{SceneGitSource, SceneHostSource}).Find(&list)
		for _, item := range list {
			if !strings.Contains(item.Content, "\"type\"") {
				continue
			}
			item.Content = strings.ReplaceAll(item.Content, "\"type\"", "\"driver\"")
			db.Select("content").Updates(item)
		}

		set.Content = `{"version":"0.8.3"}`
		db.Select("content").Updates(set)
	}

}
