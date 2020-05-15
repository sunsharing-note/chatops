package scripts

import (
	"code.rookieops.com/coolops/chatops/config"
	"code.rookieops.com/coolops/chatops/message"
	"code.rookieops.com/coolops/chatops/model"
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
	jks.ProcessMap["Job配置信息"] = jks.GetJobConfig
}

// 执行jenkins相关处理
func doJenkins(msg *message.Message) {
	// 初始化Jenkins
	initJenkins(msg)
	content := msg.ReadMessageToString()
	//var resData string
	// 对消息体进行格式化处理
	content = strings.TrimSpace(content)
	//fmt.Println(content)
	for name := range jks.ProcessMap {
		if strings.Contains(content, name) {
			switch name {
			case "build", "Job配置信息":
				_, _ = utils.Call(jks.ProcessMap, name, content)
			case "重启jenkins":
				split := strings.Split(content, " ")
				if len(split) == 2 && split[1] == "是"{
					_, _ = utils.Call(jks.ProcessMap, name)
					//rdsConn := myredis.MyPool.Get()
					//_, _ = rdsConn.Do("del", "data")
					if err := model.MyChatDao.Delete("data");err != nil{
						fmt.Println(err)
						return
					}
					if err := model.MyChatDao.Delete("name");err !=nil{
						fmt.Println(err)
						return
					}
				}else if len(split) == 2 && split[1] == "否"{
					if err := model.MyChatDao.Delete("data");err !=nil{
						fmt.Println(err)
						return
					}
					if err := model.MyChatDao.Delete("name");err !=nil{
						fmt.Println(err)
						return
					}
					//rdsConn := myredis.MyPool.Get()
					//_, _ = rdsConn.Do("DEL", "data")
				}else{
					//rdsConn := myredis.MyPool.Get()
					////_, _ = rdsConn.Do("del", "data")
					//_, _ = rdsConn.Do("set", "data", split[0])
					if err := model.MyChatDao.Set("data", split[0]);err !=nil{
						fmt.Println(err)
						return
					}
					msg.Header.Set("msgtype","text")
					msg.Body = strings.NewReader("确定要重启吗？")
					message.OutChan <- msg
				}
			default:
				_, _ = utils.Call(jks.ProcessMap, name)
			}
		}
	}
}
