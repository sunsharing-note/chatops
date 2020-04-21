package dingding

import (
	"code.rookieops.com/coolops/chatops/config"
	"code.rookieops.com/coolops/chatops/message"
	"code.rookieops.com/coolops/chatops/scripts"
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

// 定义一个钉钉结构体
type Dingtalk struct {
}

var Dingding *Dingtalk

// 加签
func signature(ts int64, secret string) string {
	strToSign := fmt.Sprintf("%d\n%s", ts, secret)
	hmac256 := hmac.New(sha256.New, []byte(secret))
	hmac256.Write([]byte(strToSign))
	data := hmac256.Sum(nil)
	return base64.StdEncoding.EncodeToString(data)
}

// 初始化dingtalk
func NewDingtalk() *Dingtalk {
	return &Dingtalk{}
}

func (d *Dingtalk) DingDing(c *gin.Context) {
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
	fmt.Println("---body/--- \r\n " + string(data))

	sign := signature(tsi, config.Setting.DingDing.AppSecret)

	// 校验成功
	if HttpSign == sign {
		//
		var body incoming
		err := json.Unmarshal(data, &body)
		if err != nil {
			fmt.Println(err)
			return
		}

		// 初始化Dingtalk
		Dingding = NewDingtalk()

		msg := message.NewMessage(body.Text.Content)
		msg.Sender = body.SenderId
		msg.Header.Set("sender", body.SenderNick)

		//content := body.Text.Content
		//fmt.Println(content)
		scripts.RunCommand(msg)
		// 开启协程，监听消息发送
		go d.listenOutChanMsg()
	}
}
