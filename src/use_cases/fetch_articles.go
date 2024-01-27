package use_cases

import (
	"enchantech-codex/src/core/database/models"
	"gorm.io/gorm"
)

type FetchArticlesUseCase struct {
	db *gorm.DB
}

func NewFetchArticlesUseCase(db *gorm.DB) *FetchArticlesUseCase {
	return &FetchArticlesUseCase{
		db: db,
	}
}

func (u *FetchArticlesUseCase) Execute(input FetchArticlesUseCaseInput) ([]models.Article, error) {
	var articles []models.Article

	articles = make([]models.Article, 0)

	err := u.db.Limit(input.Limit).Offset(input.Offset).Find(&articles).Error

	if err != nil {
		return articles, err
	}

	return articles, nil
}

type FetchArticlesUseCaseInput struct {
	// TODO add some filters
	Limit  int `schema:"limit"`
	Offset int `schema:"offset"`
}
