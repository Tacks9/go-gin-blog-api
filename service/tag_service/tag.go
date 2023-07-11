package tag_service

import (
	"encoding/json"
	"fmt"
	"go-gin-blog-api/models"
	"go-gin-blog-api/pkg/export"
	"go-gin-blog-api/pkg/gredis"
	"go-gin-blog-api/pkg/logging"
	"go-gin-blog-api/service/cache_service"
	"io"
	"strconv"
	"time"

	"github.com/tealeg/xlsx"
	"github.com/xuri/excelize/v2"
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

func (t *Tag) Export() (string, error) {
	tags, err := t.GetAll()
	if err != nil {
		return "", err
	}

	// 创建一个实例
	file := xlsx.NewFile()
	// 标签 tab
	sheet, err := file.AddSheet("标签信息")
	if err != nil {
		return "", err
	}

	// 创建一行
	row := sheet.AddRow()

	// 表头内容
	titles := []string{"ID", "标签名称", "创建人", "创建时间", "修改人", "修改时间"}

	// 创建单元格
	var cell *xlsx.Cell
	for _, title := range titles {
		// 每个单元格填写对应内容
		cell = row.AddCell()
		cell.Value = title
	}

	// 填充数据
	for _, v := range tags {
		values := []string{
			strconv.Itoa(v.ID),
			v.Name,
			v.CreatedBy,
			strconv.Itoa(v.CreatedOn),
			v.ModifiedBy,
			strconv.Itoa(v.ModifiedOn),
		}

		// 创建新的一行
		row = sheet.AddRow()
		for _, value := range values {
			// 追加新的单元格
			cell = row.AddCell()
			cell.Value = value
		}
	}

	// 获取当前时间
	time := strconv.Itoa(int(time.Now().Unix()))
	filename := "tags-" + time + ".xlsx"

	// 获取文件保存路径
	fullPath := export.GetExcelFullPath() + filename

	// 保存文件
	err = file.Save(fullPath)
	if err != nil {
		return "", err
	}

	// 返回路径
	return filename, nil
}

// 导入
func (t *Tag) Import(r io.Reader) error {
	// 读取数据流
	xlsx, err := excelize.OpenReader(r)
	if err != nil {
		return err
	}

	// 全部单元格的值
	rows, err := xlsx.GetRows("标签信息")
	if err != nil {
		logging.Info(err)
		return err
	}

	for irow, row := range rows {
		if irow > 0 {
			// 一行一行读取
			var data []string
			for _, cell := range row {
				data = append(data, cell)
			}

			// data[1] 标签名称
			// data[2] 创建人
			// models.AddTag(data[1], 1, data[2])

			// 新增标签
			t.Name = data[1]
			t.CreatedBy = data[2]
			t.State = 1

			// 判断标签是否存在
			exists, err := t.ExistByName()
			if err != nil {
				logging.Info(err)
				continue
			}
			if exists {
				logging.Info(fmt.Sprintf("%s 数据重复", t.Name))
				continue
			}
			// 数据录入
			err = t.Add()
			if err != nil {
				logging.Info(err)
				continue
			}
		}
	}

	return nil
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
