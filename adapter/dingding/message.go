package dingding

import (
	"code.rookieops.com/coolops/chatops/message"
	"fmt"
)

// At  定义需要at的用户
type At struct {
	AtMobiles []string `json:"atMobiles,omitempty"`
	IsAtAll   bool     `json:"isAtAll,omitempty"`
}

// Message 基础消息结构
type Message struct {
	MsgType string `json:"msgtype"`
	At      At     `json:"at,omitempty"`

	Text       *Text       `json:"text,omitempty"`
	Markdown   *Markdown   `json:"markdown,omitempty"`
	Link       *Link       `json:"link,omitempty"`
}

// Text text类型
type Text struct {
	Content string `json:"content,omitempty"`
}

// MarkDown 类型
type Markdown struct {
	Title string `json:"title,omitempty"`
	Text  string `json:"text,omitempty"`
}

// Link feedCard类型 links 参数
type Link struct {
	Title      string `json:"title,omitempty"`
	Text       string `json:"title,omitempty"`
	MessageURL string `json:"messageURL,omitempty"`
	PicURL     string `json:"picURL,omitempty"`
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

// NewTextMessage 新建 Text Message
func NewTextMessage(content string) *Message {
	return &Message{
		MsgType:  "text",
		At:       At{},
		Text:     &Text{Content:content},
	}
}

// NewMarkDownMessage 新建 Text Message
func NewMarkDownMessage(title,content string) *Message {
	return &Message{
		MsgType:  "markdown",
		At:       At{},
		Markdown: &Markdown{
			Title: title,
			Text:  content,
		},
	}
}

// NewLinkMessage 新建 link 类型消息
func NewLinkMessage(title, text, msgUrl, picUrl string) *Message {
	return &Message{
		MsgType: "link",
		Link: &Link{
			Title:      title,
			Text:       text,
			MessageURL: msgUrl,
			PicURL:     picUrl,
		},
	}
}

// 创建消息
func buildMessage(msg *message.Message) *Message{
	var destMsg *Message
	title := msg.Header.Get("title")
	// 如果title为空，为其制定一个默认的title
	if title == ""{
		title = "机器人消息"
	}

	// 获取msgType类型，根据类型判断选择处理方式
	msgType := msg.Header.Get("msgtype")
	fmt.Println("111111111111",msgType)
	switch  msgType {
	case "text":
		fmt.Println("发送文本消息")
		destMsg = NewTextMessage(msg.ReadMessageToString())
	case "markdown":
		fmt.Println("发送markdown消息")
		destMsg = NewMarkDownMessage(title,msg.ReadMessageToString())
	case "link":
		fmt.Println("发送link消息")
	default:
		fmt.Println("类型不匹配")
	}
	return destMsg
}