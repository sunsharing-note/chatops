package sshd

import (
	"errors"
	"fmt"
	"regexp"
)

var err error

type MySSH struct {
	SshCli   *Cli
	ShellMap map[string]interface{}
}

func NewMySSH() *MySSH {
	return &MySSH{
		ShellMap: make(map[string]interface{}),
	}
}

// GetMemoryInfo 获取内存信息
func (m *MySSH) GetMemoryInfo(content string) (output string) {
	output, err = m.SshCli.Run("free -h")
	if err != nil {
		output = err.Error()
		return
	}
	return
}

// GetDiskInfo 获取磁盘信息
func (m *MySSH) GetDiskInfo() (output string) {
	output, err = m.SshCli.Run("df -h")
	if err != nil {
		output = err.Error()
		return
	}
	return
}

// GetUpTimeInfo 获取系统负载
func (m *MySSH) GetUpTimeInfo() (output string) {
	output, err = m.SshCli.Run("uptime")
	if err != nil {
		output = err.Error()
		return
	}
	return
}

// CheckPort 检测端口是否可用
func (m *MySSH) CheckPort(content string) (output string) {
	reg := regexp.MustCompile(`\d+.\d+.\d+.\d+\s+\d+`)
	res := reg.FindAllString(content, -1)
	ipPort := regexp.MustCompile(`\s+`).Split(res[0], -1)
	fmt.Println(len(ipPort))
	if len(ipPort) == 2 {
		cmd := fmt.Sprintf("lsof -i:%s ", ipPort[1])
		_, err = m.SshCli.Run(cmd)
		if err != nil {
			output = fmt.Errorf("端口%s不存在",ipPort[1]).Error()
			return
		}
		output = fmt.Sprintf("端口%s检测通过",ipPort[1])
		return
	} else {
		output = errors.New("输入有误，请检查").Error()
	}
	return
}

