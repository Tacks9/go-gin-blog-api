package v1

import (
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

// 新增文章-表单验证
type AddArticleForm struct {
	TagID         int    `form:"tag_id" valid:"Required;Min(1)"`
	Title         string `form:"title" valid:"Required;MaxSize(100)"`
	Desc          string `form:"desc" valid:"Required;MaxSize(255)"`
	Content       string `form:"content" valid:"Required;MaxSize(65535)"`
	CreatedBy     string `form:"created_by" valid:"Required;MaxSize(100)"`
	CoverImageUrl string `form:"cover_image_url" valid:"Required;MaxSize(255)"`
	State         int    `form:"state" valid:"Range(0,1)"`
}

// 验证器-文案
var MessagesForm = map[string]string{
	"Required": "不能为空",
	"MaxSize":  "最大长度为 %d",
	"Range":    "允许的范围是:%d-%d",
}

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
	appG := app.Gin{C: c}

	// 返回数据
	data := make(map[string]interface{})

	// 实例化验证器
	validor := validation.Validation{}

	// 文章状态
	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()

		validor.Range(state, 0, 1, "state").Message("文章-状态只能为0或者1")
	}

	// 标签ID
	var tagId int = -1
	if arg := c.Query("tag_id"); arg != "" {
		tagId = com.StrTo(arg).MustInt()

		validor.Min(tagId, 1, "tag_id").Message("文章-标签ID大于0")
	}

	// 验证器-参数校验
	if validor.HasErrors() {
		app.MarkErrors(validor.Errors)
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	// 文章列表
	articleService := article_service.Article{
		TagID:    tagId,
		State:    state,
		PageNum:  util.GetPage(c),
		PageSize: setting.AppSetting.PageSize,
	}

	total, err := articleService.Count()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_COUNT_ARTICLE_FAIL, nil)
		return
	}

	articles, err := articleService.GetAll()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_ARTICLES_FAIL, nil)
		return
	}

	data["lists"] = articles
	data["total"] = total

	appG.Response(http.StatusOK, e.SUCCESS, data)

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
	var appG = app.Gin{C: c}

	// 绑定表单
	var form AddArticleForm

	// 验证器-参数校验
	httpCode, errCode, errMsg := app.BindAndValid(c, &form, MessagesForm)
	if errCode != e.SUCCESS {
		appG.FormResponse(httpCode, errCode, errMsg)
		return
	}

	// 判断是否存在
	articleService := article_service.Article{
		TagID:         form.TagID,
		Title:         form.Title,
		Desc:          form.Desc,
		Content:       form.Content,
		CoverImageUrl: form.CoverImageUrl,
		State:         form.State,
		CreatedBy:     form.CreatedBy,
	}

	// 新增
	if err := articleService.Add(); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_ARTICLE_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
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
	var appG = app.Gin{C: c}
	valid := validation.Validation{}

	id := com.StrTo(c.Param("id")).MustInt()
	tagId := com.StrTo(c.Query("tag_id")).MustInt()
	title := c.Query("title")
	desc := c.Query("desc")
	content := c.Query("content")
	modifiedBy := c.Query("modified_by")
	coverImageUrl := c.Query("cover_image_url")

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

	// 验证器
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	// 判断是否存在
	articleService := article_service.Article{
		ID:            id,
		TagID:         tagId,
		Title:         title,
		Desc:          desc,
		Content:       content,
		State:         state,
		ModifiedBy:    modifiedBy,
		CoverImageUrl: coverImageUrl,
	}
	exists, err := articleService.ExistByID()
	if err != nil {
		logging.Info(err)
		appG.Response(http.StatusInternalServerError, e.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
		return
	}
	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}

	// 编辑
	err = articleService.Edit()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EDIT_ARTICLE_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// 删除文章
// @Summary 删除文章
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/articles/{id} [delete]
func DeleteArticle(c *gin.Context) {
	appG := app.Gin{C: c}

	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("文章-ID必须大于0")

	// 验证器
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	// 判断是否存在，然后删除
	articleService := article_service.Article{ID: id}
	exists, err := articleService.ExistByID()
	if err != nil {
		logging.Info(err)
		appG.Response(http.StatusInternalServerError, e.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
		return
	}
	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}

	// 文章删除
	err = articleService.Delete()
	if err != nil {
		logging.Info(err)
		appG.Response(http.StatusInternalServerError, e.ERROR_DELETE_ARTICLE_FAIL, nil)
		return
	}

	// 数据返回
	appG.Response(http.StatusOK, e.SUCCESS, nil)

}
