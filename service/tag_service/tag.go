package tag_service

import (
	"encoding/json"
	"go-gin-blog-api/models"
	"go-gin-blog-api/pkg/gredis"
	"go-gin-blog-api/pkg/logging"
	"go-gin-blog-api/service/cache_service"
)

type Tag struct {
	ID         int
	Name       string
	CreatedBy  string
	ModifiedBy string
	State      int

	PageNum  int
	PageSize int
}

// 新增标签
func (t *Tag) Add() error {
	return models.AddTag(t.Name, t.State, t.CreatedBy)
}

// 编辑标签
func (t *Tag) Edit() error {
	// 封装编辑的数据
	data := make(map[string]interface{})
	data["modified_by"] = t.ModifiedBy
	data["name"] = t.Name
	if t.State >= 0 {
		data["state"] = t.State
	}

	// 进行修改
	return models.EditTag(t.ID, data)
}

// 删除标签
func (t *Tag) Delete() error {
	return models.DeleteTag(t.ID)
}

// 获取标签数量
func (t *Tag) Count() (int, error) {
	return models.GetTagTotal(t.getMaps())
}

// 获取标签列表
func (t *Tag) GetAll() ([]models.Tag, error) {
	var (
		tags, cacheTags []models.Tag
	)

	cache := cache_service.Tag{
		State: t.State,

		PageNum:  t.PageNum,
		PageSize: t.PageSize,
	}
	// 获取缓存
	key := cache.GetTagsKey()
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.Info(err)
		} else {
			json.Unmarshal(data, &cacheTags)
			return cacheTags, nil
		}
	}

	// 读取数据库
	tags, err := models.GetTags(t.PageNum, t.PageSize, t.getMaps())
	if err != nil {
		return nil, err
	}

	gredis.Set(key, tags, 3600)
	return tags, nil
}

// 根据 name 判断是否存在
func (t *Tag) ExistByName() (bool, error) {
	return models.ExistTagByName(t.Name)
}

// 根据 id 判断是否存在
func (t *Tag) ExistByID() (bool, error) {
	return models.ExistTagByID(t.ID)
}

// 获取Tag条件
func (t *Tag) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	maps["deleted_on"] = 0

	if t.Name != "" {
		maps["name"] = t.Name
	}
	if t.State >= 0 {
		maps["state"] = t.State
	}

	return maps
}
