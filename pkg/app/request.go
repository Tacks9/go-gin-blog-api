package app

import (
	"github.com/astaxie/beego/validation"

	"go-gin-blog-api/pkg/logging"
)

// 错误提前返回
func MarkErrors(errors []*validation.Error) {
	for _, err := range errors {
		logging.Info(err.Key, err.Message)
	}

	return
}
