package dingding

import (
	"code.rookieops.com/coolops/chatops/config"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"strconv"
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

		fmt.Println(content)
		msg := "#### 顺风耳机器人\n"+
			"> 内容：" + content
		sendMsgToDingTalk("markdown",msg)
	}
}



