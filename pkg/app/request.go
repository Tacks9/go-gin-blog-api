package app

import (
	"github.com/astaxie/beego/validation"

	"go-gin-blog-api/pkg/logging"
)

// 错误提前返回
func MarkErrors(errors []*validation.Error) {
	for _, err := range errors {
		// 报错记录
		logging.Info("err.Key :%s, err.Message:%s", err.Key, err.Message)

		// logging.Info(err.Key, err.Message)
	}

	return
}
