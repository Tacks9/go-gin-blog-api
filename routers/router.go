package routers

import (
	"go-gin-blog-api/pkg/e"
	"go-gin-blog-api/pkg/setting"
	"go-gin-blog-api/routers/api"
	v1 "go-gin-blog-api/routers/api/v1"

	"github.com/gin-gonic/gin"
)

// 初始化路由
func InitRouter() *gin.Engine {
	// 框架启动模式
	gin.SetMode(setting.RunMode)

	r := gin.New()

	// 日志中间件
	r.Use(gin.Logger())

	// 全局的恢复中间件 用于捕获异常 防止 panic 引起崩溃
	r.Use(gin.Recovery())

	// 封装测试路由
	r.GET("/test", func(ctx *gin.Context) {
		// ctx 上下文
		// 返回 json 格式化数据
		ctx.JSON(e.SUCCESS, gin.H{
			"message": e.GetMsg(e.SUCCESS) + " this is test!",
		})
	})

	// 获取 token
	r.GET("/auth", api.GetAuth)

	// 封装 API V1
	apiv1 := r.Group("/api/v1")
	{
		//标签模块
		//	获取标签列表
		apiv1.GET("/tags", v1.GetTags)
		//	新建标签
		apiv1.POST("/tags", v1.AddTag)
		//	更新指定标签
		apiv1.PUT("/tags/:id", v1.EditTag)
		//	删除指定标签
		apiv1.DELETE("/tags/:id", v1.DeleteTag)

		// 文章模块
		// 	获取文章列表
		apiv1.GET("/articles", v1.GetArticles)
		//	获取指定文章
		apiv1.GET("/articles/:id", v1.GetArticle)
		//	新建文章
		apiv1.POST("/articles", v1.AddArticle)
		//	更新指定文章
		apiv1.PUT("/articles/:id", v1.EditArticle)
		//	删除指定文章
		apiv1.DELETE("/articles/:id", v1.DeleteArticle)
	}

	return r
}
