package article_service

import (
	"encoding/json"
	"github.com/lhw0828/go-gin-example/models"
	"github.com/lhw0828/go-gin-example/pkg/gredis"
	"github.com/lhw0828/go-gin-example/pkg/logging"
	"github.com/lhw0828/go-gin-example/service/cache_service"
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

// Add 添加文章
func (a *Article) Add() error {
	article := map[string]interface{}{
		"tagID":           a.TagID,
		"title":           a.Title,
		"desc":            a.Desc,
		"content":         a.Content,
		"cover_image_url": a.CoverImageUrl,
		"state":           a.State,
		"created_by":      a.CreatedBy,
	}

	if err := models.AddArticle(article); err != nil {
		logging.Error(err)
		return err
	}
	return nil
}

func (a *Article) Edit() error {
	return models.EditArticle(a.ID, map[string]interface{}{
		"tag_id":          a.TagID,
		"title":           a.Title,
		"desc":            a.Desc,
		"state":           a.State,
		"content":         a.Content,
		"cover_image_url": a.CoverImageUrl,
		"modified_by":     a.ModifiedBy,
	})
}

func (a *Article) GetAll() ([]*models.Article, error) {
	var (
		articles, cacheArticles []*models.Article
	)

	cache := cache_service.Article{
		TagId: a.TagID,
		State: a.State,

		PageNum:  a.PageNum,
		PageSize: a.PageSize,
	}
	key := cache.GetArticlesKey()
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.Info(err)
		} else {
			json.Unmarshal(data, &cacheArticles)
			return cacheArticles, nil
		}
	}

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

// Get 获取单个文章
func (a *Article) Get() (*models.Article, error) {
	var cacheArticle *models.Article

	// 缓存
	cache := cache_service.Article{ID: a.ID}
	key := cache.GetArticleKey()

	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.Info(err)
		} else {
			err := json.Unmarshal(data, &cacheArticle)
			if err != nil {
				return nil, err
			}
			return cacheArticle, nil
		}
	}

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

// Delete 删除文章
func (a *Article) Delete() error {
	return models.DeleteArticle(a.ID)
}

// ExistByID 判断文章是否存在
func (a *Article) ExistByID() (bool, error) {
	return models.ExistArticleByID(a.ID)
}

func (a *Article) Count() (int, error) {
	return models.GetArticleTotal(a.getMaps())
}

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
