package main

import (
	. "enchantech-codex/src/core/database"
	models2 "enchantech-codex/src/core/database/models"
	"enchantech-codex/src/utils"
	"encoding/json"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"io"
)

func main() {
	err := godotenv.Load(".env", "./infra/.env")

	println("Setting up database...")
	db, err := SetupDatabase()
	utils.ErrorPanicPrinter(err, true)

	println("Seeding database...")
	err = seedPublishers(db)
	utils.ErrorPanicPrinter(err, true)
	println("Seeding database completed!")
}

func seedPublishers(db *gorm.DB) error {
	jsonFile := utils.OpenEngineerBlogs()

	byteValue, err := io.ReadAll(jsonFile)
	utils.ErrorPanicPrinter(err, true)

	var publishers []models2.Publisher
	err = json.Unmarshal(byteValue, &publishers)
	utils.ErrorPanicPrinter(err, true)

	db.Where("1 = 1").Delete(&models2.Article{})
	db.Where("1 = 1").Delete(&models2.Publisher{})

	for _, publisher := range publishers {
		result := db.Create(&publisher)

		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}
