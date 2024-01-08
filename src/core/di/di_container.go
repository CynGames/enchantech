package di

import (
	. "enchantech-codex/src/feeds/repository"
	. "enchantech-codex/src/feeds/service"
	"gorm.io/gorm"
)

type Container struct {
	FeedService *FeedService
}

func NewContainer(db *gorm.DB) *Container {
	feedRepo := NewFeedRepository(db)
	feedService := NewFeedService(feedRepo)

	return &Container{
		FeedService: feedService,
	}
}
