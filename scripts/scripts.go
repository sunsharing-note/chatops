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
	"shell":   doShell2,
	"help":    doHelpMenu,
}

func RunCommand(msg *message.Message) {
	// 查看本机磁盘/目录/文件
	//var host string
	content := msg.ReadMessageToString()
	var count int
	for menu := range menuMap {
		fmt.Println(strings.Contains(content, menu))
		if strings.Contains(content, menu) {
			_, err := utils.Call(menuMap, menu, msg)
			if err != nil {
				fmt.Println(err)
				return
			}
		}else{
			count ++
		}
	}
	if count == len(menuMap){
		doError(msg)
	}
}
