package main

import (
	"code.rookieops.com/coolops/chatops/adapter/dingding"
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

// Text text类型
type Text struct {
	Content string `json:"content,omitempty"`
}

// Body请求体
type incoming struct {
	MsgType           string            `json:"msgtype"`
	Text              *Text             `json:"text"`
	MsgId             string            `json:"msgId"`
	CreateAt          int64             `json:"createAt"`
	ConversationType  string            `json:"conversationType"` // 1-单聊、2-群聊
	ConversationId    string            `json:"conversationId"`   // // 加密的会话ID
	ConversationTitle string            `json:"conversationId"`   // 会话标题（群聊时才有）
	SenderId          string            `json:"senderId"`
	SenderNick        string            `json:"senderNick"`
	SenderCorpId      string            `json:"senderCorpId"`
	SenderStaffId     string            `json:"senderStaffId"`
	ChatbotUserId     string            `json:"chatbotUserId"`
	AtUsers           []map[string]string `json:"atUsers"`

	SessionWebhook string `json:"sessionWebhook"`
	IsAdmin        bool   `json:"isAdmin"`
}

const appSecret = "gWs8GFmPZzZmCMeRtbLKaCzor8tzJoyp4QKogD5WOVFt2Dk7UK-K5WcNgZaR3Pq2"

// 加签
func signature(ts int64, secret string) string {
	strToSign := fmt.Sprintf("%d\n%s", ts, secret)
	hmac256 := hmac.New(sha256.New, []byte(secret))
	hmac256.Write([]byte(strToSign))
	data := hmac256.Sum(nil)
	return base64.StdEncoding.EncodeToString(data)
}

func main() {
	g := gin.Default()
	g.POST("/ding/", dingtalk)
	_ = g.Run(":9999")
}

func dingtalk(c *gin.Context){
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

	sign := signature(tsi,appSecret)

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
		//go process()
		fmt.Println(content)
		dingding.SendMsgToDingtalk(content)
	}

	c.JSON(200, gin.H{
		"msgtype": "text",
		"text": `{"content": "谢谢使用此机器人}`,
	})
}