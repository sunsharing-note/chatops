package main

import (
	"code.rookieops.com/coolops/chatops/adapter"
	"code.rookieops.com/coolops/chatops/config"
	"code.rookieops.com/coolops/chatops/middleware"
	"fmt"
	"github.com/gin-gonic/gin"
)

// 程序入口，监听端口
func main() {
	fmt.Println(config.Setting.AdapterName)
	g := gin.Default()
	g.Use(middleware.LoggerToFile())
	g.POST("/ding/", process)
	_ = g.Run(":9999")
}

// 处理业务
func process(c *gin.Context) {
	adapter.Adapter(c)
}
