package v1

import (
	"fmt"
	"go-gin-blog-api/models"
	"go-gin-blog-api/pkg/app"
	"go-gin-blog-api/pkg/e"
	"go-gin-blog-api/pkg/logging"
	"go-gin-blog-api/pkg/setting"
	"go-gin-blog-api/pkg/util"
	"go-gin-blog-api/service/article_service"

	"net/http"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

// 获取单个文章
// @Summary 获取单个文章
// @Produce  json
// @Param id path  int true "ID"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/articles/{id} [get]
func GetArticle(c *gin.Context) {
	// 接收上下文
	appG := app.Gin{C: c}

	// 参数接收
	id := com.StrTo(c.Param("id")).MustInt()

	// 参数校验
	validor := validation.Validation{}
	validor.Min(id, 1, "id").Message("文章-ID必须大于0")

	// 参数有误提前返回
	if validor.HasErrors() {
		app.MarkErrors(validor.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	// 获取 Service 层
	articleService := article_service.Article{ID: id}
	exists, err := articleService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
		return
	}
	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}

	// 获取文章详情
	article, err := articleService.Get()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_ARTICLE_FAIL, nil)
		return
	}

	// 返回值
	appG.Response(http.StatusOK, e.SUCCESS, article)
}

// @Summary 获取多个文章
// @Produce  json
// @Param tag_id query int false "TagId"
// @Param state query int false "State"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/articles [get]
func GetArticles(c *gin.Context) {
	data := make(map[string]interface{})
	maps := make(map[string]interface{})

	validor := validation.Validation{}

	// 文章状态
	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		maps["state"] = state

		validor.Range(state, 0, 1, "state").Message("文章-状态只能为0或者1")
	}

	// 标签ID
	var tagId int = -1
	if arg := c.Query("tag_id"); arg != "" {
		tagId = com.StrTo(arg).MustInt()
		maps["tag_id"] = tagId

		validor.Min(tagId, 1, "tag_id").Message("文章-标签ID大于0")
	}

	// 参数校验
	code := e.INVALID_PARAMS
	if !validor.HasErrors() {
		// 获取列表
		code = e.SUCCESS
		data["lists"] = models.GetArticles(util.GetPage(c), setting.AppSetting.PageSize, maps)
		data["total"] = models.GetArticleTotal(maps)

	} else {
		for _, err := range validor.Errors {
			// log.Printf("err.Key :%s, err.Message:%s", err.Key, err.Message)
			logging.Info("err.Key :%s, err.Message:%s", err.Key, err.Message)

		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})

}

// 新增文章
// @Summary 新增文章
// @Produce  json
// @Param tag_id query int true "TagId"
// @Param title query string true "Title"
// @Param desc query string false "Desc"
// @Param content query string false "Content"
// @Param created_by query string false "CreatedBy"
// @Param state query int false "State"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/articles [post]
func AddArticle(c *gin.Context) {
	tagId := com.StrTo(c.Query("tag_id")).MustInt()
	title := c.Query("title")
	desc := c.Query("desc")
	content := c.Query("content")
	createdBy := c.Query("created_by")
	state := com.StrTo(c.DefaultQuery("state", "0")).MustInt()

	// 参数校验
	valid := validation.Validation{}
	valid.Min(tagId, 1, "tag_id").Message("文章-标签ID必须大于0")
	valid.Required(title, "title").Message("文章-标题不能为空")
	valid.Required(desc, "desc").Message("文章-简述不能为空")
	valid.Required(content, "content").Message("文章-内容不能为空")
	valid.Required(createdBy, "created_by").Message("文章-创建人不能为空")
	valid.Range(state, 0, 1, "state").Message("文章-状态只允许0或1")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		if models.ExistTagById(tagId) {
			// 封装参数
			data := make(map[string]interface{})
			data["tag_id"] = tagId
			data["title"] = title
			data["desc"] = desc
			data["content"] = content
			data["created_by"] = createdBy
			data["state"] = state

			models.AddArticle(data)
			code = e.SUCCESS
		} else {
			code = e.ERROR_NOT_EXIST_TAG
			logging.Info(fmt.Sprintf("code :%d, message:%s", code, e.GetMsg(code)))

		}
	} else {
		for _, err := range valid.Errors {
			// log.Printf("err.Key :%s, err.Message:%s", err.Key, err.Message)
			logging.Info("err.Key :%s, err.Message:%s", err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]interface{}),
	})
}

// 修改文章
// @Summary 修改文章
// @Produce  json
// @Param id path int true "ID"
// @Param tag_id query int true "TagId"
// @Param title query string true "Title"
// @Param desc query string false "Desc"
// @Param content query string false "Content"
// @Param created_by query string false "CreatedBy"
// @Param state query int false "State"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/articles/{id} [put]
func EditArticle(c *gin.Context) {
	valid := validation.Validation{}

	id := com.StrTo(c.Param("id")).MustInt()
	tagId := com.StrTo(c.Query("tag_id")).MustInt()
	title := c.Query("title")
	desc := c.Query("desc")
	content := c.Query("content")
	modifiedBy := c.Query("modified_by")

	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message("文章-状态只允许0或1")
	}

	valid.Min(id, 1, "id").Message("文章-ID必须大于0")
	valid.MaxSize(title, 100, "title").Message("文章-标题最长为100字符")
	valid.MaxSize(desc, 255, "desc").Message("文章-简述最长为255字符")
	valid.MaxSize(content, 65535, "content").Message("文章-内容最长为65535字符")
	valid.Required(modifiedBy, "modified_by").Message("文章-修改人不能为空")
	valid.MaxSize(modifiedBy, 100, "modified_by").Message("文章-修改人最长为100字符")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		if models.ExistArticleByID(id) {
			// 封装修改字段
			if models.ExistTagById(tagId) {
				data := make(map[string]interface{})
				if tagId > 0 {
					data["tag_id"] = tagId
				}
				if title != "" {
					data["title"] = title
				}
				if desc != "" {
					data["desc"] = desc
				}
				if content != "" {
					data["content"] = content
				}

				data["modified_by"] = modifiedBy

				models.EditArticle(id, data)
				code = e.SUCCESS
			} else {
				code = e.ERROR_NOT_EXIST_TAG
			}
		} else {
			code = e.ERROR_NOT_EXIST_ARTICLE
			logging.Info(fmt.Sprintf("code :%d, message:%s", code, e.GetMsg(code)))

		}
	} else {
		for _, err := range valid.Errors {
			// log.Printf("err.Key :%s, err.Message:%s", err.Key, err.Message)
			logging.Info("err.Key :%s, err.Message:%s", err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})
}

// 删除文章
// @Summary 删除文章
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/articles/{id} [delete]
func DeleteArticle(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("文章-ID必须大于0")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		// 判断是否存在，然后删除
		if models.ExistArticleByID(id) {
			models.DeleteArticle(id)
			code = e.SUCCESS
		} else {
			code = e.ERROR_NOT_EXIST_ARTICLE
			logging.Info(fmt.Sprintf("code :%d, message:%s", code, e.GetMsg(code)))
		}
	} else {
		for _, err := range valid.Errors {
			// log.Printf("err.Key :%s, err.Message:%s", err.Key, err.Message)
			logging.Info("err.Key :%s, err.Message:%s", err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})
}
