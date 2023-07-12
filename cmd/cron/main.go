package main

import (
	"go-gin-blog-api/models"
	"go-gin-blog-api/pkg/gredis"
	"go-gin-blog-api/pkg/logging"
	"go-gin-blog-api/pkg/setting"
	"time"

	"github.com/robfig/cron"
)

func main() {
	// 初始化全局配置
	setting.Setup()
	// 初始化数据库
	models.Setup()
	// 初始化日志配置
	logging.Setup()
	// Redis 设置
	gredis.Setup()

	logging.Info("Cron Starting...")

	// 创建一个定时任务
	c := cron.New()

	// 添加一个任务到 Schedule 队列中
	c.AddFunc("* * * * * *", func() {
		logging.Info("Run models.CleanAllTag...")
		models.CleanAllTag()
	})

	c.AddFunc("* * * * * *", func() {
		logging.Info("Run models.CleanAllArticle...")
		models.CleanAllArticle()
	})

	c.Start()

	// 创建一个新的定时器
	t1 := time.NewTimer(time.Second * 10)
	for {
		select {
		// 阻塞 select 等待 channel
		case <-t1.C:
			// 重置定时器，让它重新开始计时
			t1.Reset(time.Second * 10)
		}
	}

}
