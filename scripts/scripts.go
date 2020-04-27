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
			_, err := utils.Call(menuMap, menu, msg)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		//}else{
		//	fmt.Println("无效的输入，请输入help查询帮助")
		//	doError(msg)
		//	return
		//}
	}
}
