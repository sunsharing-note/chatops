package scripts

import (
	"code.rookieops.com/coolops/chatops/message"
	"code.rookieops.com/coolops/chatops/utils"
	"fmt"
	"strings"
)

// 处理脚本
var menuMap = map[string]interface{}{
	"zabbix":  doZabbix,
	"jenkins": doJenkins,
	"shell":   doShell,
	"help":    doHelpMenu,
}

func RunCommand(msg *message.Message) {
	// 查看本机磁盘/目录/文件
	//var host string
	content := msg.ReadMessageToString()
	for menu := range menuMap {
		if strings.Contains(content, menu) {
			switch menu {
			case "zabbix":
				_, _ = utils.Call(menuMap, menu, msg)
			case "jenkins":
				_, _ = utils.Call(menuMap, menu, msg)
			case "shell":
				_, _ = utils.Call(menuMap, menu, msg)
			case "help":
				_, _ = utils.Call(menuMap, menu, msg)
			default:
				fmt.Println("无效的输入，请查询help帮助")
				doError(msg)
			}

			//if err != nil {
			//	fmt.Println(err)
			//	return
			//}
		}else{
			fmt.Println("无效的输入，请输入help查询帮助")
			doError(msg)
			return
		}
	}

	//if strings.Contains(content, "help") {
	//	msg.Header.Set("msgtype", "markdown")
	//	msg.Body = strings.NewReader(menu)
	//	message.OutChan <- msg
	//
	//	// 阻塞等待用户输入
	//	// 取出消息的内容和消息的发送者进行判断
	//	// 根据用户的输入执行不同的方法
	//
	//}
	//if strings.Contains(content, "zabbix") {
	//	doZabbix(msg)
	//} else if strings.Contains(content,"jenkins"){
	//	doJenkins(msg)
	//} else {
	//	doShell(msg)
	//	fmt.Println("执行shell")
	//}
}
