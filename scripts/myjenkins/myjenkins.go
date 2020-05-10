package myjenkins

import (
	"code.rookieops.com/coolops/chatops/message"
	"fmt"
	"github.com/bndr/gojenkins"
	"strings"
	"time"
)

type MyJenkins struct {
	Jenkins    *gojenkins.Jenkins
	ProcessMap map[string]interface{}
	msg *message.Message
	output string
}

func NewMyJenkins(msg *message.Message) *MyJenkins {
	return &MyJenkins{
		ProcessMap: make(map[string]interface{}, 10),
		msg:msg,
	}
}

// 注册任务
func (m *MyJenkins) RegisterProcess(name string, value interface{}) {
	m.ProcessMap[name] = value
}

// 查询Jenkins上所有任务
func (m *MyJenkins) GetAllJob() {
	jobs, err := m.Jenkins.GetAllJobs()
	if err != nil {
		fmt.Println("获取所有Job失败，err:", err)
		m.output = "获取所有Job失败" + err.Error()
		m.sendMsgOut()
		return
	} else {
		for _, job := range jobs {
			m.output = m.output + "\n" + job.GetName() + "\n"
		}
		m.sendMsgOut()
	}
}

// 查询Jenkins上所有视图
func (m *MyJenkins) GetAllView() {
	views, err := m.Jenkins.GetAllViews()
	if err != nil {
		fmt.Println("获取所有视图失败，err:", err)
		m.output = "获取Jenkins所有视图失败"+ err.Error()
		m.sendMsgOut()
		return
	} else {
		for _, view := range views {
			m.output = m.output + "\n" + view.GetName() + "\n"
		}
		m.sendMsgOut()
	}
}

// 执行Jenkins构建工作
func (m *MyJenkins) BuildJob(content string) {
	// 获取需要build的job
	strings.TrimSpace(content)
	buildName := strings.Split(content, " ")[3]
	_, err := m.Jenkins.BuildJob(buildName)
	if err != nil {
		m.output = "build " + buildName + "失败" + err.Error()
		m.sendMsgOut()
		return
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
		m.output = buildName + "执行成功\n\n" + "详情请点击" + url
		m.sendMsgOut()
	}
}

// 重启Jenkins
func (m *MyJenkins) RestartJenkins() {
	err := m.Jenkins.SafeRestart()
	if err != nil {
		m.output = "重启Jenkins失败" + err.Error()
		m.sendMsgOut()
		return
	} else {
		m.output = "重启Jenkins成功，请稍后登录"
		m.sendMsgOut()
	}
}

// sendMsgOut  向外发送消息
func (m *MyJenkins) sendMsgOut(){
	m.msg.Header.Set("msgtype","markdown")
	m.msg.Body = strings.NewReader(m.output)
	message.OutChan <- m.msg
}