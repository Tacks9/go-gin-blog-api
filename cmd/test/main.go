package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"go-gin-blog-api/pkg/setting"
)

func main() {

	// 初始化全局配置
	setting.Setup()

	fmt.Println(setting.ServerSetting.RunMode) // debug

	r := gin.Default()

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// listen and serve on 0.0.0.0:54001
	r.Run(":54001")
}
