package main

import (
	. "enchantech-codex/src/core"
	"enchantech-codex/src/models"
	"enchantech-codex/src/utils"
	"encoding/json"
	"gorm.io/gorm"
	"io"
)

func main() {
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

	var publishers []models.Publisher
	err = json.Unmarshal(byteValue, &publishers)
	utils.ErrorPanicPrinter(err, true)

	db.Where("1 = 1").Delete(&models.Article{})
	db.Where("1 = 1").Delete(&models.Publisher{})

	for _, publisher := range publishers {
		result := db.Create(&publisher)

		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}
