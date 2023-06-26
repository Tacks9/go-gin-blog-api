package main

import (
	"fmt"
	"go-gin-blog-api/pkg/e"
	"go-gin-blog-api/pkg/setting"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// 得到一个 Router Engine
	router := gin.Default()

	// 设置路由
	router.GET("/test", func(ctx *gin.Context) {
		// ctx 上下文
		// 返回 json 格式化数据
		ctx.JSON(e.SUCCESS, gin.H{
			"message": e.GetMsg(e.SUCCESS) + " this is test!",
		})
	})

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.HTTPPort),
		Handler:        router,
		ReadTimeout:    setting.ReadTimeout,  // 允许读取的最大时间
		WriteTimeout:   setting.WriteTimeout, // 允许写入的最大时间
		MaxHeaderBytes: 1 << 20,              // 请求头的最大字节数
	}

	s.ListenAndServe()
}
