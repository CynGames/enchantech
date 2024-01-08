package di

import (
	"enchantech-codex/src/feeds/controller"
	"github.com/labstack/echo/v4"
)

func StartServer(echoInstance *echo.Echo, feedController *controller.FeedController) error {
	println("Configuring server routes...")
	echoInstance.GET("/", feedController.GetArticles)
	echoInstance.POST("/update", feedController.UpdateArticles)
	echoInstance.Static("/static", "./static")

	println("Starting server...")
	return echoInstance.Start(":11001")
}
