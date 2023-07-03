package main

import (
	"go-gin-blog-api/models"
	"go-gin-blog-api/pkg/logging"
	"time"

	"github.com/robfig/cron"
)

func main() {
	logging.Info("Starting...")

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
