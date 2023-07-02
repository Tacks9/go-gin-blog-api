package v1

import (
	"go-gin-blog-api/models"
	"go-gin-blog-api/pkg/e"
	"go-gin-blog-api/pkg/setting"
	"go-gin-blog-api/pkg/util"
	"log"
	"net/http"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

// 获取多个文章标签
func GetTags(c *gin.Context) {
	// 接收 URL Query 参数
	name := c.Query("name")

	maps := make(map[string]interface{})
	data := make(map[string]interface{})

	// 参数检查
	if name != "" {
		maps["name"] = name
	}
	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		maps["state"] = state
	}

	// 分页获取数据
	data["list"] = models.GetTags(util.GetPage(c), setting.PageSize, maps)
	data["total"] = models.GetTagTotal(maps)

	okcode := e.SUCCESS

	// 数据返回
	c.JSON(okcode, gin.H{
		"code":    okcode,
		"message": e.GetMsg(okcode),
		"data":    data,
	})
}

// 新增文章标签
func AddTag(c *gin.Context) {
	name := c.Query("name")
	state := com.StrTo(c.DefaultQuery("state", "0")).MustInt()
	createdBy := c.Query("created_by")

	// 表单验证
	validor := validation.Validation{}
	validor.Required(name, "name").Message("标签-名称不能为空")
	validor.MaxSize(name, 100, "name").Message("标签-名称最长为100字符")
	validor.Required(createdBy, "created_by").Message("标签-创建人不能为空")
	validor.MaxSize(createdBy, 100, "created_by").Message("标签-创建人最长为100字符")
	validor.Range(state, 0, 1, "state").Message("标签-状态只能为1或者0")

	code := e.INVALID_PARAMS
	if !validor.HasErrors() {
		// 存在性判定
		if !models.ExistTagByName(name) {
			code = e.SUCCESS
			models.AddTag(name, state, createdBy)
		} else {
			code = e.ERROR_EXIST_TAG
		}
	} else {
		// 参数校验
		for _, err := range validor.Errors {
			log.Println(err.Key, err.Message)
			// 数据返回
			c.JSON(http.StatusOK, gin.H{
				"code":    code,
				"message": e.GetMsg(code) + err.Key + err.Message,
				"data":    make(map[string]string),
			})
			return
		}
	}

	// 数据返回
	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": e.GetMsg(code),
		"data":    make(map[string]string),
	})
}

// 修改文章标签
func EditTag(c *gin.Context) {
	// 获取编辑的 ID
	id := com.StrTo(c.Param("id")).MustInt()

	name := c.Query("name")
	modifiedBy := c.Query("modified_by")

	// 参数校验
	validor := validation.Validation{}
	state := -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		validor.Range(state, 0, 1, "state").Message("标签-状态只能为0或者1")
	}
	validor.Required(id, "id").Message("标签-ID不能为空")
	validor.Required(modifiedBy, "modified_by").Message("标签-编辑人不能为空")
	validor.MaxSize(modifiedBy, 100, "modified_by").Message("标签-编辑人最长为100字符")
	validor.MaxSize(name, 100, "name").Message("标签-名称最长为100字符")

	code := e.INVALID_PARAMS
	if !validor.HasErrors() {
		code = e.SUCCESS
		if models.ExistTagById(id) {
			data := make(map[string]interface{})
			data["modified_by"] = modifiedBy
			if name != "" {
				data["name"] = name
			}

			if state != -1 {
				data["state"] = state
			}
			models.EditTag(id, data)
		} else {
			code = e.ERROR_NOT_EXIST_TAG
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})

}

// 删除文章标签
func DeleteTag(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	validor := validation.Validation{}
	validor.Min(id, 1, "id").Message("标签-ID必须大于0")

	code := e.INVALID_PARAMS
	if !validor.HasErrors() {
		code = e.SUCCESS
		if models.ExistTagById(id) {
			models.DeleteTag(id)
		} else {
			code = e.ERROR_NOT_EXIST_TAG
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})
}
