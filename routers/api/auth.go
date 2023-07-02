package api

import (
	"go-gin-blog-api/models"
	"go-gin-blog-api/pkg/e"
	"go-gin-blog-api/pkg/util"
	"log"
	"net/http"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

type authValid struct {
	Username string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required; MaxSize(50)"`
}

// 获取 token
func GetAuth(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	validor := validation.Validation{}

	auth := authValid{
		Username: username,
		Password: password,
	}

	ok, _ := validor.Valid(&auth)

	data := make(map[string]interface{})
	code := e.INVALID_PARAMS

	if ok {
		isExist := models.CheckAuth(username, password)
		if isExist {
			// 生成 token
			token, err := util.GenerateToken(username, password)
			if err != nil {
				code = e.ERROR_AUTH_TOKEN
				log.Println(err)

			} else {
				data["token"] = token
				code = e.SUCCESS
			}
		} else {
			code = e.ERROR_AUTH
		}
	} else {
		for _, err := range validor.Errors {
			log.Println(err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}