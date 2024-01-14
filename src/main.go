package main

import (
	database "enchantech-codex/src/core"
	scheduler "enchantech-codex/src/core"
	server "enchantech-codex/src/core"
	"enchantech-codex/src/core/di"
	"enchantech-codex/src/feeds/controller"
	"enchantech-codex/src/utils"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	err := godotenv.Load(".env", "./infra/.env")
	utils.ErrorPanicPrinter(err, false)

	db, err := database.SetupDatabase()
	utils.ErrorPanicPrinter(err, true)

	scheduler.InitializeScheduler(di.NewContainer(db).FeedService)

	container := di.NewContainer(db)

	// Mover estas instancias a di_container.go
	feedController := controller.NewFeedController(container.FeedService)
	var echoInstance = echo.New()

	err = server.StartServer(echoInstance, feedController)
	utils.ErrorPanicPrinter(err, true)

	select {}
}
