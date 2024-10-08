package tag_service

import (
	"encoding/json"
	"github.com/lhw0828/go-gin-example/models"
	"github.com/lhw0828/go-gin-example/pkg/export"
	"github.com/lhw0828/go-gin-example/pkg/gredis"
	"github.com/lhw0828/go-gin-example/pkg/logging"
	"github.com/lhw0828/go-gin-example/service/cache_service"
	"github.com/tealeg/xlsx"
	"strconv"
	"time"
)

type Tag struct {
	ID         int
	Name       string
	State      int
	CreatedBy  string
	ModifiedBy string

	PageNum  int
	PageSize int
}

func (t *Tag) ExistByName() (bool, error) {
	return models.ExistTagByName(t.Name)
}

func (t *Tag) ExistById() (bool, error) {
	return models.ExistTagById(t.ID)
}
func (t *Tag) Add() error {
	return models.AddTag(t.Name, t.State, t.CreatedBy)
}

func (t *Tag) Edit() error {
	data := make(map[string]interface{})
	data["modified_by"] = t.ModifiedBy
	data["name"] = t.Name

	if t.State >= 0 {
		data["state"] = t.State
	}
	return models.EditTag(t.ID, data)
}

func (t *Tag) Delete() error {
	return models.DeleteTag(t.ID)
}

func (t *Tag) Count() (int, error) {
	return models.GatTagTotal(t.getMaps())
}

func (t *Tag) GetAll() ([]models.Tag, error) {
	var (
		tags, cacheTags []models.Tag
	)

	cache := cache_service.Tag{
		State: t.State,

		PageNum:  t.PageNum,
		PageSize: t.PageSize,
	}

	key := cache.GetTagsKey()

	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.Info(err)
		} else {
			err := json.Unmarshal(data, &cacheTags)
			if err != nil {
				return nil, err
			}

			return cacheTags, nil
		}
	}

	tags, err := models.GetTags(t.PageNum, t.PageSize, t.getMaps())

	if err != nil {
		return nil, err
	}

	err = gredis.Set(key, tags, 3600)
	if err != nil {
		return nil, err
	}

	return tags, nil

}

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

func (t *Tag) Export() (string, error) {
	tags, err := t.GetAll()
	if err != nil {
		return "", err
	}

	file := xlsx.NewFile()
	sheet, err := file.AddSheet("标签信息")
	if err != nil {
		return "", err
	}

	titles := []string{"ID", "名称", "创建人", "创建时间", "更新人", "更新时间"}
	row := sheet.AddRow()

	var cell *xlsx.Cell
	for _, title := range titles {
		cell = row.AddCell()
		cell.Value = title
	}

	for _, v := range tags {
		values := []string{
			strconv.Itoa(int(v.ID)),
			v.Name,
			v.CreatedBy,
			strconv.Itoa(v.CreatedOn),
			v.ModifiedBy,
			strconv.Itoa(v.ModifiedOn),
		}

		row = sheet.AddRow()

		for _, value := range values {
			cell = row.AddCell()
			cell.Value = value
		}
	}
	fileTime := strconv.Itoa(int(time.Now().Unix()))
	fileName := "tag-" + fileTime + ".xlsx"

	fullPath := export.GetExcelFullPath() + fileName

	err = file.Save(fullPath)

	if err != nil {
		return "", err
	}

	return fileName, nil
}
