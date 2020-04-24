package scripts

import (
	"code.rookieops.com/coolops/chatops/config"
	"code.rookieops.com/coolops/chatops/message"
	"code.rookieops.com/coolops/chatops/scripts/myjenkins"
	"code.rookieops.com/coolops/chatops/utils"
	"fmt"
	"github.com/bndr/gojenkins"
	"reflect"
	"strings"
)

//var jenkins *gojenkins.Jenkins
var jks *myjenkins.MyJenkins

func init(){
	jks = myjenkins.NewMyJenkins()
	jks.Jenkins = gojenkins.CreateJenkins(nil, config.Setting.Jenkins.Url, config.Setting.Jenkins.UserName,config.Setting.Jenkins.PassWord)
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
	content := msg.ReadMessageToString()
	//var resData string

	for name := range jks.ProcessMap{
		if strings.Contains(content,name){
			var result []reflect.Value
			switch name {
			case "build":
				result, _ = utils.Call(jks.ProcessMap, name, content)
			default:
				result, _ = utils.Call(jks.ProcessMap, name)
			}
			msg.Header.Set("msgtype","markdown")
			msg.Body = strings.NewReader(result[0].String())
			message.OutChan <- msg
		}
	}

	//if strings.Contains(content,"所有任务"){
	//	resData = jks.GetAllJob()
	//}
	//if strings.Contains(content,"所有视图"){
	//	resData = jks.GetAllView()
	//}
	//if strings.Contains(content,"build"){
	//	strings.TrimSpace(content)
	//	buildName := strings.Split(content, " ")[3]
	//	resData = jks.BuildJob(buildName)
	//}
	//if strings.Contains(content,"重启jenkins"){
	//	resData = jks.RestartJenkins()
	//}
	//msg.Header.Set("msgtype","markdown")
	//msg.Body = strings.NewReader(resData)
	//message.OutChan <- msg
}
