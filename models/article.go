package models

import (
	"errors"
	"github.com/jinzhu/gorm"
	"time"
)

type Article struct {
	Model

	TagID int `json:"tag_id" gorm:"index"`
	Tag   Tag `json:"tag"`

	Title      string `json:"title"`
	Desc       string `json:"desc"`
	Content    string `json:"content"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      int    `json:"state"`
	DeletedOn  int    `json:"deleted_on"`
}

func (article Article) BeforeCreate(scope *gorm.Scope) error {
	err := scope.SetColumn("CreatedOn", time.Now().Unix())
	if err != nil {
		return err
	}
	return nil
}

func (article Article) BeforeUpdate(scope *gorm.Scope) error {
	err := scope.SetColumn("ModifiedOn", time.Now().Unix())
	if err != nil {
		return err
	}
	return nil
}

func ExistArticleByID(id int) (bool, error) {
	var article Article
	err := db.Select("id").Where("id = ? AND deleted_on = ?", id, 0).First(&article).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, err
	}

	if article.ID > 0 {
		return true, nil
	}
	return false, nil

}

// GetArticleTotal 根据条件获取文章总数
func GetArticleTotal(maps interface{}) (int, error) {
	var count int

	if err := db.Model(&Article{}).Where(maps).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

// GetArticles 获取文章列表
func GetArticles(pageNum int, pageSize int, maps interface{}) ([]*Article, error) {
	var articles []*Article

	err := db.Preloads("Tag").Where(maps).Offset(pageNum).Limit(pageSize).Find(&articles).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err

	}

	return articles, nil
}

// GetArticle 根据id获取文章信息
func GetArticle(id int) (*Article, error) {
	var article Article

	err := db.Where("id = ? AND deleted_on = ?", id, 0).First(&article).Error

	if err != nil {
		return nil, err
	}

	err = db.Model(&article).Related(&article.Tag).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return &article, nil
}

func EditArticle(id int, data interface{}) error {
	err := db.Model(&Article{}).Where("id = ? AND deleted_on = ?", id, 0).Updates(data).Error

	if err != nil {
		return err
	}

	return nil
}

func AddArticle(data map[string]interface{}) error {
	article := Article{
		TagID:     data["tag_id"].(int),
		Title:     data["title"].(string),
		Desc:      data["desc"].(string),
		Content:   data["content"].(string),
		CreatedBy: data["created_by"].(string),
		State:     data["state"].(int),
	}

	if err := db.Create(&article).Error; err != nil {
		return err
	}
	return nil
}

// DeleteArticle 根据id删除文章
func DeleteArticle(id int) error {
	err := db.Where("id = ? AND deleted_on = ?", id, 0).Delete(&Article{}).Error

	if err != nil {
		return err
	}

	return nil
}

func CleanAllArticle() error {
	err := db.Unscoped().Where("deleted_on != ?", 0).Delete(&Article{}).Error

	if err != nil {
		return err
	}

	return nil
}
