package sshd

import (
	"code.rookieops.com/coolops/chatops/message"
	"fmt"
	"github.com/relex/aini"
	"net"
	"regexp"
	"strings"
	"time"
)

var err error

type MySSH struct {
	SshCli   *Cli
	ShellMap map[string]interface{}
	File *aini.InventoryData
	msg *message.Message
	output string
}

func NewMySSH(msg *message.Message) *MySSH {
	return &MySSH{
		ShellMap: make(map[string]interface{}),
		msg:msg,
	}
}

// GetMemoryInfo 获取内存信息
func (m *MySSH) GetMemoryInfo(content string) {
	ipReg := regexp.MustCompile(`((25[0-5]|2[0-4]\d|[01]?\d\d?)\.){3}(25[0-5]|2[0-4]\d|[01]?\d\d?)`)
	ipList := ipReg.FindAllString(content,-1)
	ipRes := m.checkIP(ipList)
	for _,ip:= range ipRes {
		m.login(ip)
		m.output, err = m.SshCli.Run("free -h")
		if err != nil {
			m.output = err.Error()
			m.sendMsgOut()
			return
		}
		m.sendMsgOut()
	}
}

// GetDiskInfo 获取磁盘信息
func (m *MySSH) GetDiskInfo(content string) {
	ipReg := regexp.MustCompile(`((25[0-5]|2[0-4]\d|[01]?\d\d?)\.){3}(25[0-5]|2[0-4]\d|[01]?\d\d?)`)
	ipList := ipReg.FindAllString(content,-1)
	ipRes := m.checkIP(ipList)
	for _,ip:= range ipRes {
		m.login(ip)
		m.output, err = m.SshCli.Run("df -h")
		if err != nil {
			m.output = err.Error()
			m.sendMsgOut()
			return
		}
		m.sendMsgOut()
	}
}

// login 登录服务器
func (m *MySSH) login(ip string){
	// 获取用户名，密码，端口
	host := m.File.Match(ip)[0]
	sshHost := host.Name
	sshPort := host.Port
	sshUser := host.Vars["ssh_user"]
	sshPassword := host.Vars["ssh_password"]
	// 连接服务器
	address := fmt.Sprintf("%s:%d", sshHost, sshPort)
	m.SshCli = NewSSH(sshUser, sshPassword, address)
}

// GetUpTimeInfo 获取系统负载
func (m *MySSH) GetUpTimeInfo(content string) {
	ipReg := regexp.MustCompile(`((25[0-5]|2[0-4]\d|[01]?\d\d?)\.){3}(25[0-5]|2[0-4]\d|[01]?\d\d?)`)
	ipList := ipReg.FindAllString(content,-1)
	ipRes := m.checkIP(ipList)
	for _,ip:= range ipRes{
		m.login(ip)
		m.output, err = m.SshCli.Run("uptime")
		if err != nil {
			m.output = err.Error()
			m.sendMsgOut()
		}
		m.sendMsgOut()
	}
}

// sendMsgOut
func (m *MySSH)sendMsgOut(){
	tmp := "顺风耳机器人\n" +
		"输出内容：\n" + m.output
	m.msg.Header.Set("msgtype", "text")
	m.msg.Body = strings.NewReader(tmp)
	message.OutChan <- m.msg
}

// checkIP 检测IP是否可用
func (m *MySSH)checkIP(ipList []string) (res []string) {
	// 检测IP是否在配置文件中，如果不在则返回无该IP，并从数组中删除该IP
	for index, ip := range ipList {
		hosts := m.File.Match(ip)
		if len(hosts) == 0 {
			m.output = fmt.Sprintf("无效的主机IP" + ip + ",请检查。")
			m.sendMsgOut()
			res = append(ipList[:index], ipList[index+1:]...)
		}else{
			res = ipList
		}
	}
	return
}

// CheckPort 检测端口是否可用
func (m *MySSH) CheckPort(content string)  {
	var (
		ipPort []string
		conn net.Conn
	)
	reg := regexp.MustCompile(`\d+.\d+.\d+.\d+:\d+`)
	res := reg.FindAllString(content, -1)
	for _, ip := range res{
		//ipPort = regexp.MustCompile(`\s+`).Split(ip, -1)
		//checkIp := fmt.Sprintf("%s:%s",ipPort[0],ipPort[1])
		ipPort = strings.Split(ip,":")
		conn, err = net.DialTimeout("tcp", ip, 3*time.Second);
		if err !=nil{
			// 表示未被占用
			m.output = fmt.Sprintf("IP为%s的服务器%s未被使用",ipPort[0],ipPort[1])
			m.sendMsgOut()
			return
		}else{
			m.output = fmt.Sprintf("IP为%s的服务器%s已经被使用",ipPort[0],ipPort[1])
			m.sendMsgOut()
		}
		_ = conn.Close()
	}
}

// CheckServer 检测服务是否可用
func (m *MySSH) CheckServer(content string){
	ipReg := regexp.MustCompile(`((25[0-5]|2[0-4]\d|[01]?\d\d?)\.){3}(25[0-5]|2[0-4]\d|[01]?\d\d?)`)
	domainReg := regexp.MustCompile(`[a-zA-Z0-9][a-zA-Z0-9_-]{0,62}(\.[a-zA-Z0-9][a-zA-Z0-9_-]{0,62})*(\.[a-zA-Z][a-zA-Z0-9]{0,10}){1}`)
	ipList := ipReg.FindAllString(content,-1)
	domainList := domainReg.FindAllString(content,-1)
	if len(ipList) != 0{
		for _, ip := range ipList{
			conn, err := net.Dial("ip:icmp", ip)
			if err != nil {
				// 检测不通过
				m.output = fmt.Sprintf("IP: %s 检测结果：异常",ip)
				m.sendMsgOut()
				return
			}else{
				// 检测通过
				m.output = fmt.Sprintf("IP: %s 检测结果：正常",ip)
				m.sendMsgOut()
			}
			_ = conn.Close()
		}
	}
	if len(domainList) != 0{
		for _,domain := range domainList{
			conn, err := net.Dial("ip:icmp", domain)
			if err != nil {
				// 检测不通过
				m.output = fmt.Sprintf("域名：%s 检测结果：异常",domain)
				m.sendMsgOut()
				return
			}else{
				// 检测通过
				m.output = fmt.Sprintf("域名：%s 对应IP：%s 检测结果：正常 ",domain,conn.RemoteAddr().String())
				m.sendMsgOut()
			}
			_ = conn.Close()
		}
	}
}