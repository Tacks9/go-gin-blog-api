package jwt

import (
	"go-gin-blog-api/pkg/e"
	"go-gin-blog-api/pkg/util"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// jwt 中间件
func JWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var code int
		var data interface{}

		code = e.SUCCESS

		// 获取 token
		token := ctx.Query("token")
		if token == "" {
			code = e.INVALID_PARAMS
		} else {
			claims, err := util.ParseToken(token)
			if err != nil {
				code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
			} else if time.Now().Unix() > claims.ExpiresAt {
				code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
			}
		}

		// 认证失败
		if code != e.SUCCESS {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": code,
				"msg":  e.GetMsg(code),
				"data": data,
			})

			// 中断
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
