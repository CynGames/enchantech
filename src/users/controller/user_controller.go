package controller

import (
	"enchantech-codex/src/users/service"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

//func (fc *FeedController) GetArticles(context echo.Context) error {
//	println("Fetching articles...")
//	articles, err := fc.feedService.GetArticles()
//
//	if err != nil {
//		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
//	}
//
//	return templates.MainPage(articles).Render(context.Request().Context(), context.Response())
//}
//
//func (fc *FeedController) UpdateArticles(context echo.Context) error {
//	println("Updating articles...")
//	err := fc.feedService.GetRSSXMLContent()
//
//	if err != nil {
//		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
//	}
//
//	return context.NoContent(http.StatusOK)
//}
