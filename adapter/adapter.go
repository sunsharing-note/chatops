package adapter

import (
	"code.rookieops.com/coolops/chatops/adapter/dingding"
	"code.rookieops.com/coolops/chatops/config"
	"fmt"
	"github.com/gin-gonic/gin"
)

// 根据配置文件配置的adapter做不同的处理
func Adapter(c *gin.Context){
	if config.Setting.AdapterName == "dingding"{
		dingding.Dingding.DingDing(c)

	}else if config.Setting.AdapterName == "wechat"{
		fmt.Println("处理企业微信")
	}
}
