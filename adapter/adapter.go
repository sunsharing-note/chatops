package adapter

import (
	"code.rookieops.com/coolops/chatops/adapter/dingding"
	"code.rookieops.com/coolops/chatops/config"
	"code.rookieops.com/coolops/chatops/middleware"
	"fmt"
	"github.com/gin-gonic/gin"
)

// 根据配置文件配置的adapter做不同的处理
func Adapter(c *gin.Context) {
	switch config.Setting.AdapterName {
	case "dingding":
		dingding.Dingding.DingDing(c)
	case "wechat":
		fmt.Println("处理企业微信")
	default:
		// 不处理任务逻辑，可以优雅返回一个参数
		middleware.Logger().Error("adapter not exists")
	}
}
