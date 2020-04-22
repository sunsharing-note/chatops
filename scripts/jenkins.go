package scripts

import (
	"code.rookieops.com/coolops/chatops/config"
	"code.rookieops.com/coolops/chatops/message"
	"fmt"
	"github.com/bndr/gojenkins"
	"strings"
)

var jenkins *gojenkins.Jenkins

func init(){
	jenkins = gojenkins.CreateJenkins(nil, config.Setting.Jenkins.Url, config.Setting.Jenkins.UserName,config.Setting.Jenkins.PassWord)
	// 初始化
	_, err := jenkins.Init()
	if err != nil {
		fmt.Println("init jenkins failed.")
		return
	}
}

func doJenkins(msg *message.Message) {
	content := msg.ReadMessageToString()
	var resData string
	var tmp string
	if strings.Contains(content,"所有job"){
		jobs, err := jenkins.GetAllJobs()
		if err != nil {
			fmt.Println("获取所有Job失败，err:",err)
			resData = "获取所有Job失败"
		}else{
			for _,job := range jobs{
				resData = resData+"\n"+job.GetName()+"\n"
			}
		}
		tmp = "#### 顺风耳机器人\n" +
			//"##### 主机：" + ip + "\n" +
			"所有Job:\n" + resData
	}
	if strings.Contains(content,"所有视图"){
		views, err := jenkins.GetAllViews()
		if err != nil {
			fmt.Println("获取所有视图失败，err:",err)
			resData = "获取所有视图失败"
		}else{
			for _,view := range views{
				resData = resData+"\n"+view.GetName()+"\n"
			}
		}
		tmp = "#### 顺风耳机器人\n" +
			//"##### 主机：" + ip + "\n" +
			"所有视图:\n" + resData
	}
	if strings.Contains(content,"build"){
		// 获取需要build的job
		strings.TrimSpace(content)
		buildName := strings.Split(content, " ")[3]
		_, err := jenkins.BuildJob(buildName)
		if err != nil {
			resData = "build "+buildName+"失败"
		}else{
			resData = "成功"
		}
		tmp = "#### 顺风耳机器人\n" +
			//"##### 主机：" + ip + "\n" +
			"build:\n" + resData
	}
	msg.Header.Set("msgtype","markdown")
	msg.Body = strings.NewReader(tmp)
	message.OutChan <- msg
}
