package article_service

import (
	"encoding/json"
	"go-gin-blog-api/models"
	"go-gin-blog-api/pkg/gredis"
	"go-gin-blog-api/pkg/logging"
	"go-gin-blog-api/service/cache_service"
)

type Article struct {
	ID            int
	TagID         int
	Title         string
	Desc          string
	Content       string
	CoverImageUrl string
	State         int
	CreatedBy     string
	ModifiedBy    string

	PageNum  int
	PageSize int
}

// Add 新增文章
func (a *Article) Add() error {
	article := map[string]interface{}{
		"tag_id":          a.TagID,
		"title":           a.Title,
		"desc":            a.Desc,
		"content":         a.Content,
		"created_by":      a.CreatedBy,
		"cover_image_url": a.CoverImageUrl,
		"state":           a.State,
	}

	// 添加数据
	if err := models.AddArticle(article); err != nil {
		return err
	}

	return nil
}

// Edit 编辑文章
func (a *Article) Edit() error {
	return models.EditArticle(a.ID, map[string]interface{}{
		"tag_id":          a.TagID,
		"title":           a.Title,
		"desc":            a.Desc,
		"content":         a.Content,
		"cover_image_url": a.CoverImageUrl,
		"state":           a.State,
		"modified_by":     a.ModifiedBy,
	})

}

// Get 获取一篇文章
func (a *Article) Get() (*models.Article, error) {
	var cacheArticle *models.Article

	// 获取缓存 KEY
	cache := cache_service.Article{ID: a.ID}
	key := cache.GetArticleKey()
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.Info(err)
		} else {
			// 存在直接返回
			json.Unmarshal(data, &cacheArticle)
			return cacheArticle, nil
		}
	}

	// 不存在读取数据库
	article, err := models.GetArticle(a.ID)
	if err != nil {
		return nil, err
	}

	err = gredis.Set(key, article, 3600)
	if err != nil {
		return nil, err
	}

	return article, nil
}

// GetAll() 获取一批文章
func (a *Article) GetAll() ([]*models.Article, error) {
	var cacheArticles []*models.Article

	// 获取缓存 KEY
	cache := cache_service.Article{
		TagID: a.TagID,
		State: a.State,

		PageNum:  a.PageNum,
		PageSize: a.PageSize,
	}
	key := cache.GetArticlesKey()

	// 判断是否存在
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.Info(err)
		} else {
			// 存在直接返回
			json.Unmarshal(data, &cacheArticles)
			return cacheArticles, nil
		}
	}

	// 不存在读取数据库
	articles, err := models.GetArticles(a.PageNum, a.PageSize, a.getMaps())
	if err != nil {
		return nil, err
	}

	err = gredis.Set(key, articles, 3600)
	if err != nil {
		return nil, err
	}

	return articles, nil
}

// Delete 删除
func (a *Article) Delete() error {
	return models.DeleteArticle(a.ID)
}

// ExistByID 判断是否存在
func (a *Article) ExistByID() (bool, error) {
	return models.ExistArticleByID(a.ID)
}

// Count 获取数量
func (a *Article) Count() (int, error) {
	return models.GetArticleTotal(a.getMaps())
}

// 获取检查条件
func (a *Article) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	maps["deleted_on"] = 0
	if a.State != -1 {
		maps["state"] = a.State
	}
	if a.TagID != -1 {
		maps["tag_id"] = a.TagID
	}

	return maps
}
