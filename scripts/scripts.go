package scripts

import (
	"code.rookieops.com/coolops/chatops/message"
	"fmt"
	"strings"
)

// 处理脚本

func RunCommand(msg *message.Message) {
	// 查看本机磁盘/目录/文件
	//var host string
	menu := `## 帮助信息
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
- 查询jenkins的所有job
- 查询jenkins的所有视图
- 查询jenkins视图[view_name]下的所有job
- 执行jenkins build [job_name]
- 重启jenkins
------------------
请按着帮助信息输入内容！
`
	content := msg.ReadMessageToString()
	if strings.Contains(content, "help") {
		msg.Header.Set("msgtype", "markdown")
		msg.Body = strings.NewReader(menu)
		message.OutChan <- msg

		// 阻塞等待用户输入
		// 取出消息的内容和消息的发送者进行判断
		// 根据用户的输入执行不同的方法

	}
	if strings.Contains(content, "zabbix") {
		doZabbix(msg)
	} else if strings.Contains(content,"jenkins"){
		doJenkins(msg)
	} else {
		doShell(msg)
		fmt.Println("执行shell")
	}
}
