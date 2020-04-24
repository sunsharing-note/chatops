package myjenkins

import (
	"fmt"
	"github.com/bndr/gojenkins"
	"strings"
	"time"
)

type MyJenkins struct {
	Jenkins    *gojenkins.Jenkins
	ProcessMap map[string]interface{}
}

func NewMyJenkins() *MyJenkins {
	return &MyJenkins{
		ProcessMap: make(map[string]interface{}, 10),
	}
}

// 注册任务
func (m *MyJenkins) RegisterProcess(name string, value interface{}) {
	m.ProcessMap[name] = value
}

// 查询Jenkins上所有任务
func (m *MyJenkins) GetAllJob() (resData string) {
	jobs, err := m.Jenkins.GetAllJobs()
	if err != nil {
		fmt.Println("获取所有Job失败，err:", err)
		resData = "获取所有Job失败"
	} else {
		for _, job := range jobs {
			resData = resData + "\n" + job.GetName() + "\n"
		}
	}
	resData = "#### 顺风耳机器人\n" +
		//"##### 主机：" + ip + "\n" +
		"所有任务列表:\n" + resData
	return
}

// 查询Jenkins上所有视图
func (m *MyJenkins) GetAllView() (resData string) {
	views, err := m.Jenkins.GetAllViews()
	if err != nil {
		fmt.Println("获取所有视图失败，err:", err)
		resData = "获取所有视图失败"
	} else {
		for _, view := range views {
			resData = resData + "\n" + view.GetName() + "\n"
		}
	}
	resData = "#### 顺风耳机器人\n" +
		//"##### 主机：" + ip + "\n" +
		"所有视图:\n" + resData
	return
}

// 执行Jenkins构建工作
func (m *MyJenkins) BuildJob(content string) (resData string) {
	// 获取需要build的job
	strings.TrimSpace(content)
	buildName := strings.Split(content, " ")[3]
	_, err := m.Jenkins.BuildJob(buildName)
	if err != nil {
		resData = "build " + buildName + "失败"
	} else {
		// 获取build的url
		_, _ = m.Jenkins.Poll()
		time.Sleep(time.Second * 10)
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

// 重启Jenkins
func (m *MyJenkins) RestartJenkins() (resData string) {
	err := m.Jenkins.SafeRestart()
	if err != nil {
		resData = "重启Jenkins失败"
	} else {
		resData = "重启Jenkins成功，请稍后登录"
	}
	resData = "#### 顺风耳机器人\n" + resData
	return
}