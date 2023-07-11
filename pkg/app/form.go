package app

import (
	"go-gin-blog-api/pkg/e"
	"net/http"
	"reflect"
	"strings"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

// 绑定参数，并且验证
func BindAndValid(c *gin.Context, form interface{}, messages map[string]string) (int, int, string) {

	var errMsg string

	// 绑定表单
	err := c.Bind(form)
	if err != nil {
		return http.StatusBadRequest, e.INVALID_PARAMS, errMsg
	}

	// 设置错误提示
	validation.SetDefaultMessage(messages)

	// 表单验证
	valid := validation.Validation{}

	check, err := valid.Valid(form)
	if err != nil {
		return http.StatusInternalServerError, e.ERROR, errMsg
	}

	// 表单验证-有误
	if !check && len(valid.Errors) > 0 {
		MarkErrors(valid.Errors)

		// 获取报错的字段
		fieldName := valid.Errors[0].Field
		fieldMsg := valid.Errors[0].Message

		// 获取结构体
		structForm := reflect.TypeOf(form)

		// 获取真正报错的字段
		field, _ := structForm.Elem().FieldByName(fieldName)
		fieldTag := field.Tag.Get("form")

		// errMsg = fmt.Sprintf("%s %s", fieldTag, fieldMsg)

		// 重设报错信息
		errMsg = strings.Replace(fieldMsg, fieldName, fieldTag, 1)

		return http.StatusBadRequest, e.INVALID_PARAMS, errMsg
	}

	return http.StatusOK, e.SUCCESS, errMsg
}
