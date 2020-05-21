package message

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/textproto"
	"strings"
)

// 定义消息结构体
type Message struct {
	From       string    // 消息来源
	To         string    // 消息接收者
	Sender     string    //	消息发送者
	Header     Header    // 消息头
	KeepHeader bool      // 如果为True，消息得Header在一次会话结束之前不会清除
	Body       io.Reader // 消息体
}

var OutChan = make(chan *Message,20)
var InputChan = make(chan *Message,10)
var Msg *Message
// Header 消息附带的头信息，键-值对
type Header map[string][]string

// NewMessage 初始化消息
func NewMessage(content string,to ...string)*Message{
	 msg := &Message{
		 From:       "",
		 To:         "",
		 Sender:     "",
		 Header:     Header{},
		 KeepHeader: false,
		 Body:       strings.NewReader(content),
	 }
	 if len(to)>0{
	 	msg.To = to[0]
	 }
	 return msg
}

// 读取消息内容
func (m *Message) ReadMessageToString()string{
	content, err := ioutil.ReadAll(m.Body)
	if err != nil {
		fmt.Println(err)
	}
	m.Body = bytes.NewBuffer(content)

	return string(content)
}

// Get 从头信息中获取与给定键关联的第一个值
func (h Header) Get(key string) string {
	return textproto.MIMEHeader(h).Get(key)
}

// Set 将key设置为单个值，它替换与key的现有值
func (h Header) Set(key, value string) {
	textproto.MIMEHeader(h).Set(key, value)
}

func SendMsg(msg *Message,msgType,info string){
	msg.Header.Set("msgtype",msgType)
	msg.Body = strings.NewReader(info)
	OutChan <- msg
}