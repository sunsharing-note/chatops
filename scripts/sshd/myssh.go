package sshd

type MySSH struct {
	SshCli *Cli
	ShellMap map[string]interface{}
}

func NewMySSH()*MySSH{
	return &MySSH{
		ShellMap: make(map[string]interface{}),
	}
}

// GetMemoryInfo 获取内存信息
func (m *MySSH)GetMemoryInfo(content string)(output string){
	output, err := m.SshCli.Run("free -h")
	if err != nil {
		output = err.Error()
		return
	}
	return
}

// GetDiskInfo 获取磁盘信息
func (m *MySSH)GetDiskInfo()(output string){
	output, err:= m.SshCli.Run("df -h")
	if err != nil {
		output = err.Error()
		return
	}
	return
}

// GetUpTimeInfo 获取系统负载
func (m *MySSH)GetUpTimeInfo()(output string){
	output, err:= m.SshCli.Run("uptime")
	if err != nil {
		output = err.Error()
		return
	}
	return
}

