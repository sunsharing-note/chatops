package scripts

import (
	"code.rookieops.com/coolops/chatops/message"
	"code.rookieops.com/coolops/chatops/scripts/sshd"
	"fmt"
	"regexp"
	"strings"
)

func doShell(msg *message.Message){
	content := msg.ReadMessageToString()
	var command string
	if strings.Contains(content,"磁盘信息"){
		command = "df -h"
	}
	if strings.Contains(content,"内存信息"){
		command = "free -h"
	}
	// 找到主机IP
	reg := regexp.MustCompile(`\d+.\d+.\d+.\d+`)
	res := reg.FindAllString(content,-1)
	for _,ip:=range res {
		fmt.Println(ip)
		address := fmt.Sprintf("%s:%s", ip, "22")
		cli := sshd.NewSSH("root", "coolops@123456", address)
		output, err := cli.Run(command)
		if err != nil {
			content = "执行命令失败"
		}
		tmp := "顺风耳机器人\n" +
			"查询主机：" + ip + "\n" +
			"输出内容：\n" +
			output
		//dingding.SendMsgToDingTalk("markdown", msg)
		// 直接放进channel中
		msg.Header.Set("msgtype","text")
		msg.Body = strings.NewReader(tmp)
		message.OutChan <- msg
	}
}
