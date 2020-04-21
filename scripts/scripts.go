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
	content := msg.ReadMessageToString()
	if strings.Contains(content,"zabbix"){
		doZabbix(msg)
	}else{
		doShell(msg)
		fmt.Println("执行shell")
	}
}