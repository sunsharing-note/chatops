package main

import (
	"code.rookieops.com/coolops/chatops/adapter"
	"code.rookieops.com/coolops/chatops/config"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println(config.Setting.AdapterName)
	g := gin.Default()
	g.POST("/ding/", process)
	_ = g.Run(":9999")
}

// 处理业务
func process(c *gin.Context){
	adapter.Adapter(c)
}
