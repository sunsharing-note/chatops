package scripts

import (
	"code.rookieops.com/coolops/chatops/config"
	"code.rookieops.com/coolops/chatops/message"
	"code.rookieops.com/coolops/chatops/scripts/sshd"
	"code.rookieops.com/coolops/chatops/utils"
	"fmt"
	"github.com/relex/aini"
	"strings"
)

var (
	mySSH *sshd.MySSH
)

func initSSH(msg *message.Message) {
	// 初始化ssh
	mySSH = sshd.NewMySSH(msg)
	mySSH.ShellMap["内存信息"] = mySSH.GetMemoryInfo
	mySSH.ShellMap["磁盘信息"] = mySSH.GetDiskInfo
	mySSH.ShellMap["负载信息"] = mySSH.GetUpTimeInfo
	mySSH.ShellMap["端口检测"] = mySSH.CheckPort
	mySSH.ShellMap["服务检测"] = mySSH.CheckServer
	// 加载主机配置文件
	loadHostsFile()
}

// 加载主机配置文件
func loadHostsFile() {
	var err error
	mySSH.File, err = aini.ParseFile(config.Setting.SSH.FilePath)
	if err != nil {
		fmt.Println(err)
		return
	}
}


func doShell(msg *message.Message){
	// 初始化
	initSSH(msg)
	content := msg.ReadMessageToString()
	//var outputMsg string
	// 对消息体进行格式化处理
	content = strings.TrimSpace(content)
	for name := range mySSH.ShellMap {
		// 判断关键字是否存在
		if strings.Contains(content, name) {
			switch name {
			case "端口检测" , "服务检测":
				_, _ = utils.Call(mySSH.ShellMap, name, content)
				//outputMsg = output[0].String()
			case "内存信息" , "磁盘信息", "负载信息":
				// 获取IP -> 登录服务器 -> 执行命令 -> 返回结果
				_, _ = utils.Call(mySSH.ShellMap, name, content)
				//outputMsg = output[0].String()
			default:

			}
		}
	}
}