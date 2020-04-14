package dingding

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