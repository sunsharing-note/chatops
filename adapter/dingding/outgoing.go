package dingding

import (
	"net/http"
	"strings"
)

// 向钉钉发消息

//var baseHookUrl = "https://oapi.dingtalk.com/robot/send?access_token=61f0415bbdd8c05317a086a63b042c154ca22ddee6ffd0915d67c20e9040e1ae"

func sendMsgToDingTalk(msg string){
	//请求地址模板
	webHook := `https://oapi.dingtalk.com/robot/send?access_token=61f0415bbdd8c05317a086a63b042c154ca22ddee6ffd0915d67c20e9040e1ae`
	content := `{"msgtype": "text",
		"text": {"content": "`+ msg + `"}
	}`
	//创建一个请求
	req, err := http.NewRequest("POST", webHook, strings.NewReader(content))
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
