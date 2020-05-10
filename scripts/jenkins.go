package scripts

import (
	"code.rookieops.com/coolops/chatops/config"
	"code.rookieops.com/coolops/chatops/message"
	"code.rookieops.com/coolops/chatops/scripts/myjenkins"
	"code.rookieops.com/coolops/chatops/utils"
	"fmt"
	"github.com/bndr/gojenkins"
	"strings"
)

var jks *myjenkins.MyJenkins

func initJenkins(msg *message.Message) {
	jks = myjenkins.NewMyJenkins(msg)
	jks.Jenkins = gojenkins.CreateJenkins(nil, config.Setting.Jenkins.Url, config.Setting.Jenkins.UserName, config.Setting.Jenkins.PassWord)
	// 初始化
	_, err := jks.Jenkins.Init()
	if err != nil {
		fmt.Println("init myjenkins failed.")
		return
	}
	// 注册任务
	jks.ProcessMap["所有任务"] = jks.GetAllJob
	jks.ProcessMap["所有视图"] = jks.GetAllView
	jks.ProcessMap["build"] = jks.BuildJob
	jks.ProcessMap["重启jenkins"] = jks.RestartJenkins
}

// 执行jenkins相关处理
func doJenkins(msg *message.Message) {
	// 初始化Jenkins
	initJenkins(msg)
	content := msg.ReadMessageToString()
	//var resData string

	for name := range jks.ProcessMap {
		if strings.Contains(content, name) {
			switch name {
			case "build":
				_, _ = utils.Call(jks.ProcessMap, name, content)
			default:
				_, _ = utils.Call(jks.ProcessMap, name)
			}
		}
	}
}
