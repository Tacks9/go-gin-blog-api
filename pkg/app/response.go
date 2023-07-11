package app

import (
	"github.com/gin-gonic/gin"

	"go-gin-blog-api/pkg/e"
)

type Gin struct {
	C *gin.Context
}

// 返回格式
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// 输出 JSON
// Response setting gin.JSON
func (g *Gin) Response(httpCode, errCode int, data interface{}) {
	g.C.JSON(httpCode, gin.H{
		"code": errCode,
		"msg":  e.GetMsg(errCode),
		"data": data,
	})

	return
}

// 输出 JSON
func (g *Gin) FormResponse(httpCode, errCode int, errMsg string) {
	g.C.JSON(httpCode, gin.H{
		"code": errCode,
		"msg":  e.GetMsg(errCode) + "-" + errMsg,
		"data": nil,
	})

	return
}
