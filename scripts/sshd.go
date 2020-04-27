package scripts

import (
	"code.rookieops.com/coolops/chatops/message"
	"code.rookieops.com/coolops/chatops/scripts/sshd"
	"fmt"
	"regexp"
	"strings"
)
var shellMenu = map[string]string{
	"内存信息":"free ",
	"磁盘信息":"df ",
	"进程信息":"ps -ef | grep ",
}

func runShell(ip,content,command string)(tmp string){
	// 根据IP到数据库中查找端口，用户名，密码
	address := fmt.Sprintf("%s:%s", ip, "22")
	cli := sshd.NewSSH("root", "coolops@123456", address)
	output, err := cli.Run(command)
	if err != nil {
		content = "执行命令失败"
		fmt.Println(err)
		return
	}
}

// 检查IP是否有效
func checkIP(ipList []string,msg *message.Message) (res []string) {
	// 检测IP是否在配置文件中，如果不在则返回无该IP，并从数组中删除该IP
	for index,ip := range ipList{
		hosts := file.Match(ip)
		if len(hosts) == 0{
			tmp := "无效的主机IP"+ip+",请检查。"
			msg.Header.Set("msgtype", "text")
			msg.Body = strings.NewReader(tmp)
			message.OutChan <- msg
			res = append(ipList[:index],ipList[index+1:]...)
		}
	}
	return
}

func doShell(msg *message.Message){
	content := msg.ReadMessageToString()
	reg := regexp.MustCompile(`\d+.\d+.\d+.\d+`)
	res := reg.FindAllString(content,-1)
	for menu := range shellMenu{
		if strings.Contains(content,menu){
			for _,ip:=range res{
				tmp := runShell(ip,content,shellMenu[menu])
				msg.Header.Set("msgtype","text")
				msg.Body = strings.NewReader(tmp)
				message.OutChan <- msg
			}
		}
	}
}
