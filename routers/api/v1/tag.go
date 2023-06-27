package v1

import (
	"go-gin-blog-api/models"
	"go-gin-blog-api/pkg/e"
	"go-gin-blog-api/pkg/setting"
	"go-gin-blog-api/pkg/util"

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
}

// 修改文章标签
func EditTag(c *gin.Context) {
}

// 删除文章标签
func DeleteTag(c *gin.Context) {
}
