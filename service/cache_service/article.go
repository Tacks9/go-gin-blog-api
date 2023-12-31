package cache_service

import (
	"go-gin-blog-api/pkg/e"
	"strconv"
	"strings"
)

type Article struct {
	ID    int
	TagID int
	State int

	PageNum  int
	PageSize int
}

// [缓存KEY] 获取单个文章 Key
func (a *Article) GetArticleKey() string {
	return e.CACHE_ARTICLE + "_" + strconv.Itoa(a.ID)
}

// [缓存KEY] 获取文章列表 Key
func (a *Article) GetArticlesKey() string {
	keys := []string{
		e.CACHE_ARTICLE,
		"LIST",
	}

	if a.ID > 0 {
		keys = append(keys, strconv.Itoa(a.ID))
	}
	if a.TagID > 0 {
		keys = append(keys, strconv.Itoa(a.TagID))
	}
	if a.State >= 0 {
		keys = append(keys, strconv.Itoa(a.State))
	}
	if a.PageNum > 0 {
		keys = append(keys, strconv.Itoa(a.PageNum))
	}
	if a.PageSize > 0 {
		keys = append(keys, strconv.Itoa(a.PageSize))
	}

	return strings.Join(keys, "_")
}
