package main

import (
	database "enchantech-codex/src/core"
	scheduler "enchantech-codex/src/core"
	server "enchantech-codex/src/core"
	"enchantech-codex/src/core/di"
	"enchantech-codex/src/feeds/controller"
	"enchantech-codex/src/utils"
	"github.com/labstack/echo/v4"
)

func main() {
	db, err := database.SetupDatabase()
	utils.ErrorPanicPrinter(err, true)

	scheduler.InitializeScheduler(di.NewContainer(db).FeedService)

	container := di.NewContainer(db)
	feedController := controller.NewFeedController(container.FeedService)

	var echoInstance = echo.New()
	err = server.StartServer(echoInstance, feedController)
	utils.ErrorPanicPrinter(err, true)

	select {}
}
