package di

import (
	. "enchantech-codex/src/feeds/repository"
	. "enchantech-codex/src/feeds/service"
	"gorm.io/gorm"
)

type Container struct {
	FeedService *FeedService
}

// retornar todo (controller, service, repository) en un GET-X

func NewContainer(db *gorm.DB) *Container {
	feedRepo := NewFeedRepository(db)
	feedService := NewFeedService(feedRepo)

	return &Container{
		FeedService: feedService,
	}
}
