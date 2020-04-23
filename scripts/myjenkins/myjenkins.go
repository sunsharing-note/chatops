package myjenkins

import (
	"fmt"
	"github.com/bndr/gojenkins"
	"time"
)
type MyJenkins struct {
	Jenkins *gojenkins.Jenkins
}

func NewMyJenkins()*MyJenkins{
	return &MyJenkins{
	}
}

func (m *MyJenkins)GetAllJob()(resData string){
	jobs, err := m.Jenkins.GetAllJobs()
	if err != nil {
		fmt.Println("获取所有Job失败，err:",err)
		resData = "获取所有Job失败"
	}else{
		for _,job := range jobs{
			resData = resData+"\n"+job.GetName()+"\n"
		}
	}
	resData = "#### 顺风耳机器人\n" +
		//"##### 主机：" + ip + "\n" +
		"所有任务列表:\n" + resData
	return
}

func (m *MyJenkins)GetAllView()(resData string){
	views, err := m.Jenkins.GetAllViews()
	if err != nil {
		fmt.Println("获取所有视图失败，err:",err)
		resData = "获取所有视图失败"
	}else{
		for _,view := range views{
			resData = resData+"\n"+view.GetName()+"\n"
		}
	}
	resData = "#### 顺风耳机器人\n" +
		//"##### 主机：" + ip + "\n" +
		"所有视图:\n" + resData
	return
}

func (m *MyJenkins)BuildJob(buildName string)(resData string){
	// 获取需要build的job

	_, err := m.Jenkins.BuildJob(buildName)
	if err != nil {
		resData = "build "+buildName+"失败"
	}else{
		// 获取build的url
		_, _ = m.Jenkins.Poll()
		time.Sleep(time.Second*10)
		job, err := m.Jenkins.GetJob(buildName)
		if err != nil {
			fmt.Println("获取job失败")
			return
		}
		_, _ = job.Poll()
		build, err := job.GetLastBuild()
		if err != nil {
			fmt.Println("获取最后一次Build失败")
			return
		}
		_, _ = build.Poll()
		url := build.GetUrl()
		resData = buildName + "执行成功\n\n" + "详情请点击" + url
	}
	resData = "#### 顺风耳机器人>\n" +
		//"##### 主机：" + ip + "\n" +
		"build:\n" + resData
	return
}

func (m *MyJenkins)RestartJenkins()(resData string){
	err := m.Jenkins.SafeRestart()
	if err != nil {
		resData = "重启Jenkins失败"
	}else{
		resData = "重启Jenkins成功，请稍后登录"
	}
	resData = "#### 顺风耳机器人\n" + resData
	return
}