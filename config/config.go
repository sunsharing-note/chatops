package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Adapter
	DingDing
	WeChat
}

type Adapter struct {
	AdapterName string `yaml:"adapter_name"`
}

type DingDing struct {
	AppSecret string `yaml:"app_secret"`
	AccessToken string `yaml:"access_token"`
}

type WeChat struct {
	AppSecret string `yaml:"app_secret"`
}

// 定义一个全局变量
var Setting Config

// 初始化配置
func init(){

	file, err := ioutil.ReadFile("./config/chatops.yaml")
	if err != nil {
		fmt.Println("open config file failed. err:",err)
		return
	}
	_ = yaml.Unmarshal(file, &Setting)
}