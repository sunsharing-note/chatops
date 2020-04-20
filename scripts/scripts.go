package scripts

import "strings"

// 处理脚本

func RunCommand(content string){
	// 查看本机磁盘/目录/文件
	//var host string
	if strings.Contains(content,"监控"){
		doZabbix(content)
	}else{
		doShell(content)
	}

}