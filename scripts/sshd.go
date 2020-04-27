package scripts

import (
	"code.rookieops.com/coolops/chatops/config"
	"code.rookieops.com/coolops/chatops/message"
	"code.rookieops.com/coolops/chatops/scripts/sshd"
	"code.rookieops.com/coolops/chatops/utils"
	"fmt"
	"github.com/relex/aini"
	"regexp"
	"strings"
)

var (
	mySSH *sshd.MySSH
	file  *aini.InventoryData
)

func init() {
	// 加载主机配置文件
	loadHostsFile()
	// 初始化ssh
	mySSH = sshd.NewMySSH()
	mySSH.ShellMap["内存信息"] = mySSH.GetMemoryInfo
	mySSH.ShellMap["磁盘信息"] = mySSH.GetDiskInfo
	mySSH.ShellMap["负载信息"] = mySSH.GetUpTimeInfo
	mySSH.ShellMap["端口检测"] = mySSH.CheckPort
}

// 加载主机配置文件
func loadHostsFile() {
	var err error
	file, err = aini.ParseFile(config.Setting.SSH.FilePath)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// 检查IP是否有效
func checkIP(ipList []string, msg *message.Message) (res []string) {
	// 检测IP是否在配置文件中，如果不在则返回无该IP，并从数组中删除该IP
	for index, ip := range ipList {
		hosts := file.Match(ip)
		if len(hosts) == 0 {
			tmp := "无效的主机IP" + ip + ",请检查。"
			msg.Header.Set("msgtype", "text")
			msg.Body = strings.NewReader(tmp)
			message.OutChan <- msg
			res = append(ipList[:index], ipList[index+1:]...)
		}else{
			res = ipList
		}
	}
	return
}

func runShell(ip, name, content string) []string {
	// 根据IP到数据库中查找端口，用户名，密码
	resData := make([]string, 0)
	// 检测IP是否存在于配置中
	hosts := file.Match(ip)
	fmt.Println("222222222",name)
	for _, host := range hosts {
		// 获取端口，用户名，密码
		sshHost := host.Name
		sshPort := host.Port
		sshUser := host.Vars["ssh_user"]
		sshPassword := host.Vars["ssh_password"]
		fmt.Println(sshPort, sshUser, sshPassword)
		// 连接服务器
		address := fmt.Sprintf("%s:%d", sshHost, sshPort)
		mySSH.SshCli = sshd.NewSSH(sshUser, sshPassword, address)
		var msg string
		switch name {
		case "端口检测":
			output, _ := utils.Call(mySSH.ShellMap, name, content)
			msg = output[0].String()
		default:
			output, _ := utils.Call(mySSH.ShellMap, name)
			msg = output[0].String()
		}
		//fmt.Println(msg)
		resData = append(resData, msg)
	}
	return resData
}

func doShell(msg *message.Message) {
	content := msg.ReadMessageToString()
	// 获取content中的IP地址
	reg := regexp.MustCompile(`\d+.\d+.\d+.\d+`)
	res := reg.FindAllString(content, -1)
	ipList := checkIP(res, msg)
	for name := range mySSH.ShellMap {
		// 判断关键字是否存在
		if strings.Contains(content, name) {
			for _, ip := range ipList {
				resData := runShell(ip, name, content)
				for _, tmp := range resData {
					tmp = "顺风耳机器人\n" +
						"查询主机：" + ip + "\n" +
						"输出内容：\n" + tmp
					msg.Header.Set("msgtype", "text")
					msg.Body = strings.NewReader(tmp)
					message.OutChan <- msg
				}
			}
		}
	}
}
