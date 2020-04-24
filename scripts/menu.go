package scripts

import (
	"code.rookieops.com/coolops/chatops/message"
	"strings"
)

var menu = `## 帮助信息
------------------
### 1、zabbix
- 获取zabbix的所有用户
- 获取zabbix的版本信息
- 获取zabbix[IP]的主机信息
- 获取zabbix[IP]的报警信息
- 获取zabbix[IP]的事件信息
- 获取zabbix[IP]的历史记录
### 2、处理Kubernetes
### 3、执行Linux命令
- 获取[IP]的内存信息
- 获取[IP]的磁盘信息
### 4、处理Jenkins
- 查询jenkins的所有任务
- 查询jenkins的所有视图
- 查询jenkins视图[view_name]下的所有任务
- 执行jenkins build [job_name]
- 重启jenkins
------------------
请按着帮助信息输入内容！
`

func doMenu(msg *message.Message){
	msg.Header.Set("msgtype", "markdown")
	msg.Body = strings.NewReader(menu)
	message.OutChan <- msg
}
