package dingding

import (
	"bytes"
	"code.rookieops.com/coolops/chatops/config"
	"code.rookieops.com/coolops/chatops/message"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// 向钉钉发消息

var baseHookUrl = "https://oapi.dingtalk.com/robot/send"

// 定义监听outChan，有消息就发送到群中
func (d *Dingtalk) listenOutChanMsg(){
	for {
		select {
		case out := <-message.OutChan:
			fmt.Println(out)
			d.SendMsgToDingTalk(out)
		default:

		}
	}
}

func (d *Dingtalk)SendMsgToDingTalk(outMsg *message.Message){
	//请求地址模板
	accessToekn := config.Setting.DingDing.AccessToken
	query := url.Values{}
	query.Set("access_token",accessToekn)
	hookUrl, _ := url.Parse(baseHookUrl)
	hookUrl.RawQuery = query.Encode()
	message := buildMessage(outMsg)
	msgContent,_ := json.Marshal(message)
	//创建一个请求
	req, err := http.NewRequest("POST", hookUrl.String(), bytes.NewReader(msgContent))
	if err != nil {
		// handle error
	}

	client := &http.Client{}
	//设置请求头
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	//发送请求
	resp, err := client.Do(req)
	//关闭请求
	defer resp.Body.Close()

	if err != nil {
		// handle error
	}
}
