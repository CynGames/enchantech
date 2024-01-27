package main

import (
	scheduler "enchantech-codex/src/core"
	server "enchantech-codex/src/core"
	database2 "enchantech-codex/src/core/database"
	"enchantech-codex/src/core/di"
	"enchantech-codex/src/utils"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	utils.ErrorPanicPrinter(err, false)

	db, err := database2.SetupDatabase()
	utils.ErrorPanicPrinter(err, true)
	scheduler.InitializeScheduler(di.NewContainer(db).FeedService)
	container := di.NewContainer(db)

	err = server.StartServer(container)
	utils.ErrorPanicPrinter(err, true)

	select {}
}
