package main

import (
	"fmt"
	"go-gin-blog-api/pkg/setting"
	"go-gin-blog-api/routers"
	"net/http"
)

func main() {
	// 获取路由
	router := routers.InitRouter()

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.HTTPPort),
		Handler:        router,
		ReadTimeout:    setting.ReadTimeout,  // 允许读取的最大时间
		WriteTimeout:   setting.WriteTimeout, // 允许写入的最大时间
		MaxHeaderBytes: 1 << 20,              // 请求头的最大字节数
	}

	s.ListenAndServe()
}
