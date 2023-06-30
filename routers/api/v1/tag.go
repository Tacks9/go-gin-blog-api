package v1

import (
	"go-gin-blog-api/models"
	"go-gin-blog-api/pkg/e"
	"go-gin-blog-api/pkg/setting"
	"go-gin-blog-api/pkg/util"

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
	createBy := c.Query("create_by")

	validor := validation.Validation{}
	validor.Required(name, "name").Message("标签-名称不能为空")
	validor.MaxSize(name, 100, "name").Message("标签-名称最长为100字符")
	validor.Required(createBy, "create_by").Message("标签-创建人不能为空")
	validor.MaxSize(createBy, 100, "create_by").Message("标签-创建人最长为100字符")
	validor.Range(state, 0, 1, "state").Message("标签-状态只能为1或者0")

	code := e.INVALID_PARAMS
	if !validor.HasErrors() {
		if !models.ExistTagByName(name) {
			code = e.SUCCESS
			models.AddTag(name, state, createBy)
		} else {
			code = e.ERROR_EXIST_TAG
		}
	}

	// 数据返回
	c.JSON(code, gin.H{
		"code":    code,
		"message": e.GetMsg(code),
		"data":    make(map[string]string),
	})
}

// 修改文章标签
func EditTag(c *gin.Context) {
}

// 删除文章标签
func DeleteTag(c *gin.Context) {
}
