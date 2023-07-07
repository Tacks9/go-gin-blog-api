package cache_service

import (
	"go-gin-blog-api/pkg/e"
	"strconv"
	"strings"
)

type Tag struct {
	ID    int
	Name  string
	State int

	PageNum  int
	PageSize int
}

// 获取标签的 Key
func (t *Tag) GetTagsKey() string {
	keys := []string{
		e.CACHE_TAG,
		"LIST",
	}

	if t.Name != "" {
		keys = append(keys, t.Name)
	}

	if t.State >= 0 {
		keys = append(keys, strconv.Itoa(t.State))
	}

	if t.PageNum > 0 {
		keys = append(keys, strconv.Itoa(t.PageNum))
	}
	if t.PageSize > 0 {
		keys = append(keys, strconv.Itoa(t.PageSize))
	}

	// 数组转成字符串
	return strings.Join(keys, "_")

}
