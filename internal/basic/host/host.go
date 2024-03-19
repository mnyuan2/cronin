package host

import (
	"cron/internal/basic/errs"
	"golang.org/x/crypto/ssh"
)

type Config struct {
	Ip     string `json:"ip"`
	Port   string `json:"port"`
	User   string `json:"user"`
	Secret string `json:"secret"`
}

type Host struct {
	conf *Config
}

func NewHost(conf *Config) *Host {
	return &Host{
		conf: conf,
	}
}

// 远程执行
func (m *Host) RemoteExec(statement string) ([]byte, errs.Errs) {
	config := &ssh.ClientConfig{
		User: m.conf.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(m.conf.Secret),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	// 连接到远程服务器
	conn, er := ssh.Dial("tcp", m.conf.Ip+":"+m.conf.Port, config)
	if er != nil {
		return nil, errs.New(er, "拨号失败")
	}
	defer conn.Close()

	// 创建一个新的会话
	session, er := conn.NewSession()
	if er != nil {
		return nil, errs.New(er, "创建会话失败")
	}
	defer session.Close()

	// 执行Shell脚本
	output, er := session.CombinedOutput(statement)
	if er != nil {
		return nil, errs.New(er, "执行脚本失败")
	}
	// 打印脚本输出
	return output, nil
}
