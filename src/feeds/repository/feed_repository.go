package repository

import (
	models2 "enchantech-codex/src/core/database/models"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type FeedRepository struct {
	db *gorm.DB
}

func NewFeedRepository(db *gorm.DB) *FeedRepository {
	return &FeedRepository{db: db}
}

func (fr *FeedRepository) GetArticles() ([]models2.Article, error) {
	var articles []models2.Article
	err := fr.db.Find(&articles).Error

	return articles, err
}

func (fr *FeedRepository) GetPublishers() ([]models2.Publisher, error) {
	var publishers []models2.Publisher
	err := fr.db.Find(&publishers).Error

	return publishers, err
}

func (fr *FeedRepository) CreateArticles(articles []models2.Article) error {
	for _, article := range articles {
		result := fr.db.Create(&article)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {

				continue
			} else {
				return fmt.Errorf("error creating article: %w", result.Error)
			}
		}
	}
	return nil
}

func (fr *FeedRepository) CreatePublisher(publishers []models2.Publisher) error {
	result := fr.db.Create(&publishers)
	return result.Error
}

func (fr *FeedRepository) RemoveAll() {
	fr.db.Where("1 = 1").Delete(&models2.Publisher{})
	fr.db.Where("1 = 1").Delete(&models2.Article{})
}
