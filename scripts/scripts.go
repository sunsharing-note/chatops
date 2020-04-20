package scripts

import "strings"

// 处理脚本

func RunCommand(content string)(msg []string){
	// 查看本机磁盘/目录/文件
	//var host string
	if strings.Contains(content,"zabbix"){
		msg = doZabbix(content)
	}else{
		msg = doShell(content)
	}
	return msg
}