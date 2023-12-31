package main

import (
	"fmt"
	"go-gin-blog-api/models"
	"go-gin-blog-api/pkg/gredis"
	"go-gin-blog-api/pkg/logging"
	"go-gin-blog-api/pkg/setting"
	"go-gin-blog-api/routers"
	"net/http"
)

// @title go-gin-blog-api
// @version 1.0
// @description An example blog of gin
// @securityDefinitions.apikey ApiKeyAuth
// @in query
// @name token
// @contact.name Tacks9
// @contact.url https://github.com/Tacks9/go-gin-blog-api
func main() {

	// 初始化全局配置
	setting.Setup()
	// 初始化数据库
	models.Setup()
	// 初始化日志配置
	logging.Setup()
	// Redis 设置
	gredis.Setup()

	logging.Info("Api Starting...")

	// 获取路由
	router := routers.InitRouter()

	// 优雅重启
	// fvbock/endless 热更新是采取创建子进程后，将原进程退出的方式
	// endless.DefaultReadTimeOut = setting.ServerSetting.ReadTimeout
	// endless.DefaultWriteTimeOut = setting.ServerSetting.WriteTimeout
	// endless.DefaultMaxHeaderBytes = 1 << 20
	// endPoint := fmt.Sprintf(":%d", setting.ServerSetting.HttpPort)
	// // 实例化对象 endlessServer
	// server := endless.NewServer(endPoint, router)

	// // 输出当前进程 PID
	// server.BeforeBegin = func(add string) {
	// 	log.Printf("Actual pid is %d", syscall.Getpid())
	// }

	// // 监听 HTTP
	// err := server.ListenAndServe()
	// if err != nil {
	// 	log.Printf("Server err: %v", err)
	// }

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.ServerSetting.HttpPort),
		Handler:        router,
		ReadTimeout:    setting.ServerSetting.ReadTimeout,  // 允许读取的最大时间
		WriteTimeout:   setting.ServerSetting.WriteTimeout, // 允许写入的最大时间
		MaxHeaderBytes: 1 << 20,                            // 请求头的最大字节数
	}

	s.ListenAndServe()

	// 优雅关闭
	// go func() {
	// 	if err := s.ListenAndServe(); err != nil {
	// 		log.Printf("Listen: %s\n", err)
	// 	}
	// }()
	// // 创建 chan 接收操作系统中断信号
	// quit := make(chan os.Signal)
	// // 将接收到的中断信号发送到 quit 通道
	// signal.Notify(quit, os.Interrupt)
	// // 程序阻塞，直到收到 中断信号
	// <-quit

	// // 准备关闭服务器
	// log.Println("Shutdown Server ...")

	// // 等待所有正在处理的请求完成后再关闭服务器，超过5s报错
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()
	// if err := s.Shutdown(ctx); err != nil {
	// 	log.Fatal("Server Shutdown:", err)
	// }
	// log.Println("Server exiting")

}
