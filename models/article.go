package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Article struct {
	Model

	// gorm 索引
	TagID int `json:"tag_id" gorm:"index"`

	// 模型关联 类名+ID 的方式去找到这两个类之间的关联关系
	Tag Tag `json:"tag"`

	Title      string `json:"title"`
	Desc       string `json:"desc"`
	Content    string `json:"content"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      int    `json:"state"`
}

// 获取单个文章 根据ID
func GetArticle(id int) (article Article) {
	db.Where("id = ?", id).First(&article)
	// 获取关联的标签表数据
	db.Model(&article).Related(&article.Tag)

	return
}

// 创建文章
func AddArticle(data map[string]interface{}) bool {
	db.Create(&Article{
		// 类型断言 V.(I)
		TagID:     data["tag_id"].(int),
		Title:     data["title"].(string),
		Desc:      data["desc"].(string),
		Content:   data["content"].(string),
		CreatedBy: data["created_by"].(string),
		State:     data["state"].(int),
	})

	return true
}

// 编辑某个文章
func EditArticle(id int, data interface{}) bool {
	db.Model(&Article{}).Where("id = ?", id).Updates(data)

	return true
}

// 删除某个文章
func DeleteArticle(id int) bool {
	db.Where("id = ?", id).Delete(Article{})

	return true
}

// 获取文章的数量
func GetArticleTotal(maps interface{}) (count int) {
	db.Model(&Article{}).Where(maps).Count(&count)

	return
}

// 获取文章的列表-分页
func GetArticles(pageNum int, pageSize int, maps interface{}) (articles []Article) {
	// Preload 预加载器，执行两次 select 对 article/tag
	db.Preload("Tag").Where(maps).Offset(pageNum).Limit(pageSize).Find(&articles)

	return
}

// 判断文章是否存在 根据 ID
func ExistArticleByID(id int) bool {
	var article Article
	db.Select("id").Where("id = ?", id).First(&article)

	if article.ID > 0 {
		return true
	}

	return false
}

// 创建之前
func (article *Article) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("CreatedOn", time.Now().Unix())

	return nil
}

// 更新之前
func (article *Article) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("ModifiedOn", time.Now().Unix())

	return nil
}
