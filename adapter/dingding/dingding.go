package dingding

import (
	"code.rookieops.com/coolops/chatops/config"
	"code.rookieops.com/coolops/chatops/scripts/sshd"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"
)

// 加签
func signature(ts int64, secret string) string {
	strToSign := fmt.Sprintf("%d\n%s", ts, secret)
	hmac256 := hmac.New(sha256.New, []byte(secret))
	hmac256.Write([]byte(strToSign))
	data := hmac256.Sum(nil)
	return base64.StdEncoding.EncodeToString(data)
}

func DingDing(c *gin.Context){
	// 获取body里的请求参数
	//fmt.Println(c.Request.Header)
	HttpSign := c.Request.Header.Get("Sign")
	HttpTimestamp := c.Request.Header.Get("Timestamp")

	// timestamp 与系统当前时间戳如果相差1小时以上，则认为是非法的请求。
	tsi, err := strconv.ParseInt(HttpTimestamp, 10, 64)
	if err != nil {
		log.Printf("请求头可能未附加时间戳信息!!")
	}

	data, _ := ioutil.ReadAll(c.Request.Body)
	fmt.Println("---body/--- \r\n "+string(data))

	sign := signature(tsi,config.Setting.DingDing.AppSecret)

	// 校验成功
	if HttpSign == sign{
		//
		var body incoming
		err := json.Unmarshal(data, &body)
		if err != nil {
			fmt.Println(err)
			return
		}
		content := body.Text.Content
		// 起一个协程去执行任务
		//shellCommand := []string{
		//	"df",
		//	"ls",
		//	"cat",
		//}
		fmt.Println(content)
		// 查看本机磁盘/目录/文件
		//var host string
		var command string
		if strings.Contains(content,"磁盘信息"){
			command = "df -h"
		}
		if strings.Contains(content,"内存信息"){
			command = "free -h"
		}
		// 找到主机IP
		reg := regexp.MustCompile(`\d+.\d+.\d+.\d`)
		res := reg.FindAllString(content,-1)
		for _,ip:=range res{
			fmt.Println(ip)
			address := fmt.Sprintf("%s:%s",ip,"22")
			cli := sshd.NewSSH("root","coolops@123456",address)
			output, err := cli.Run(command)
			if err != nil {
				content = "执行命令失败"
			}
			fmt.Println(output)
			content = output
			msg := "#### 顺风耳机器人\n"+
				"##### 主机：" + ip +"\n"+
				"> 内容：" + content
			sendMsgToDingTalk("markdown",msg)
		}


	}
}



