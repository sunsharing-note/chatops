package scripts

import (
	"code.rookieops.com/coolops/chatops/message"
	"strings"
)

//var helpMenu = map[string]string

var helpMenu = `## 帮助信息
------------------
### 1、zabbix
- 获取zabbix的所有用户
- 获取zabbix的版本信息
- 获取zabbix[IP]的主机信息
- 获取zabbix[IP]的报警信息
- 获取zabbix[IP]的事件信息
- 获取zabbix[IP]的历史记录
------------------
### 2、处理Kubernetes
------------------
### 3、执行Linux命令
- shell 内存信息 [IP]
- shell 磁盘信息 [IP]
- shell 端口检测 [IP|DOMAIN]:[PORT]
- shell 日志信息 [IP] [LOGPATH]
- shell 负载信息 [IP]
- shell 服务检测 [IP|DOMAIN]
------------------
### 4、处理Jenkins
- 查询jenkins的所有任务
- 查询jenkins的所有视图
- 查询jenkins视图[view_name]下的所有任务
- 构建:jenkins build [job_name]
- jenkins Job配置信息,Name:[job_name]
- 重启jenkins
------------------
请按着帮助信息输入内容！
`

func doHelpMenu(msg *message.Message){
	msg.Header.Set("msgtype", "markdown")
	msg.Body = strings.NewReader(helpMenu)
	message.OutChan <- msg
}

func doError(msg *message.Message){
	tmp := "无效的输入，请查询help帮助"
	msg.Header.Set("msgtype", "text")
	msg.Body = strings.NewReader(tmp)
	message.OutChan <- msg
}
