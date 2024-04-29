package biz

import (
	"bytes"
	"cron/internal/basic/enum"
	"cron/internal/pb"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"testing"
)

func TestSqlParse(t *testing.T) {
	list1 := []*pb.CronConfigListItem{}

	str := `[
    {
        "id": 123,
        "name": "git sql test",
        "spec": "",
        "type": 5,
        "remark": "",
        "status": 1,
        "command": {
            "cmd": "",
            "rpc": {
                "addr": "",
                "body": "",
                "proto": "",
                "action": "",
                "header": [],
                "method": "GRPC",
                "actions": []
            },
            "sql": {
                "driver": "mysql",
                "source": {
                    "id": 2,
                    "port": "",
                    "title": "",
                    "database": "",
                    "hostname": "",
                    "password": "",
                    "username": ""
                },
                "interval": 0,
                "statement": [],
                "err_action": 1,
                "statement_git": [
                    {
                        "ref": "master",
                        "path": [
                            "2023/sm1201/test.sql",
                            "2023/sm1201/test2.sql",
                            ""
                        ],
                        "owner": "zhubaoe",
                        "link_id": 12,
                        "project": "russell"
                    }
                ],
                "err_action_name": "终止任务",
                "statement_source": "git"
            },
            "http": {
                "url": "",
                "body": "",
                "header": [
                    {
                        "key": "",
                        "value": ""
                    }
                ],
                "method": "GET"
            },
            "jenkins": {
                "name": "",
                "params": [
                    {
                        "key": "",
                        "value": ""
                    }
                ],
                "source": {
                    "id": 0
                }
            }
        },
        "msg_set": [],
        "entry_id": 0,
        "protocol": 4,
        "status_dt": "",
        "type_name": "模块",
        "update_dt": "2024-03-01 22:42:55",
        "top_number": 0,
        "status_name": "停用",
        "protocol_name": "sql",
        "status_remark": "",
        "top_error_number": 0
    },
    {
        "id": 121,
        "name": "kobe-order-test",
        "spec": "",
        "type": 5,
        "remark": "",
        "status": 1,
        "command": {
            "cmd": "",
            "rpc": {
                "addr": "",
                "body": "",
                "proto": "",
                "action": "",
                "header": [],
                "method": "GRPC",
                "actions": []
            },
            "sql": {
                "driver": "mysql",
                "source": {
                    "id": 0,
                    "port": "",
                    "title": "",
                    "database": "",
                    "hostname": "",
                    "password": "",
                    "username": ""
                },
                "interval": 0,
                "statement": [],
                "err_action": 1,
                "statement_git": [],
                "err_action_name": "",
                "statement_source": "local"
            },
            "http": {
                "url": "",
                "body": "",
                "header": [
                    {
                        "key": "",
                        "value": ""
                    }
                ],
                "method": "GET"
            },
            "jenkins": {
                "name": "kobe-service-common",
                "params": [
                    {
                        "key": "BRANCH",
                        "value": "feature/test"
                    },
                    {
                        "key": "ENV",
                        "value": "dev"
                    },
                    {
                        "key": "SERVICENAME",
                        "value": "order"
                    },
                    {
                        "key": "",
                        "value": ""
                    }
                ],
                "source": {
                    "id": 11
                }
            }
        },
        "msg_set": [
            {
                "msg_id": 8,
                "status": 1,
                "notify_user_ids": []
            }
        ],
        "entry_id": 0,
        "protocol": 5,
        "status_dt": "",
        "type_name": "模块",
        "update_dt": "2024-02-27 10:56:33",
        "top_number": 5,
        "status_name": "停用",
        "protocol_name": "jenkins",
        "status_remark": "",
        "top_error_number": 4
    }
]`

	if err := jsoniter.Unmarshal([]byte(str), &list1); err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Println(list1)

	sql := `select 1; 
select 2; select 3;`

	list := [][]byte{}
	if 1 == enum.BoolYes {
		list = bytes.Split([]byte(sql), []byte(";"))
	} else {
		list = [][]byte{[]byte(sql)}
	}
	for i, item := range list {
		s := bytes.TrimSpace(item)
		if s != nil {
			fmt.Println(i, s, string(s))
		}
	}
	fmt.Println(list)
}
