package dingding

import (
	"bytes"
	"code.rookieops.com/coolops/chatops/config"
	"encoding/json"
	"net/http"
	"net/url"
)

// 向钉钉发消息

var baseHookUrl = "https://oapi.dingtalk.com/robot/send"

func SendMsgToDingTalk(title,msg string){
	//请求地址模板
	accessToekn := config.Setting.DingDing.AccessToken
	query := url.Values{}
	query.Set("access_token",accessToekn)
	hookUrl, _ := url.Parse(baseHookUrl)
	hookUrl.RawQuery = query.Encode()
	message := buildMessage(title, msg)
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
