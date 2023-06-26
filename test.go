package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// listen and serve on 0.0.0.0:54001
	r.Run(":54001")
}
