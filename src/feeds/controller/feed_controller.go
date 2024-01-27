package controller

import (
	"enchantech-codex/src/feeds/service"
	"github.com/labstack/echo/v4"
	"net/http"
)

type FeedController struct {
	feedService *service.FeedService
}

func NewFeedController(feedService *service.FeedService) *FeedController {
	return &FeedController{
		feedService: feedService,
	}
}

func (fc *FeedController) UpdateArticles(context echo.Context) error {
	println("Updating articles...")
	err := fc.feedService.GetRSSXMLContent()

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return context.NoContent(http.StatusOK)
}
