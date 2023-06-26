package routers

import (
	"go-gin-blog-api/pkg/e"
	"go-gin-blog-api/pkg/setting"

	"github.com/gin-gonic/gin"
)

// 初始化路由
func InitRouter() *gin.Engine {
	r := gin.New()

	// 日志中间件
	r.Use(gin.Logger())

	// 全局的恢复中间件 用于捕获异常 防止 panic 引起崩溃
	r.Use(gin.Recovery())

	// 框架启动模式
	gin.SetMode(setting.RunMode)

	// 封装测试路由
	r.GET("/test", func(ctx *gin.Context) {
		// ctx 上下文
		// 返回 json 格式化数据
		ctx.JSON(e.SUCCESS, gin.H{
			"message": e.GetMsg(e.SUCCESS) + " this is test!",
		})
	})

	return r
}
