package database

import (
	models2 "enchantech-codex/src/core/database/models"
	"enchantech-codex/src/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

func SetupDatabase() (*gorm.DB, error) {
	var parsedAddress string
	var err error

	if os.Getenv("IS_PROD") == "true" {
		dbURI := os.Getenv("DATABASE_URI")
		parsedAddress, err = utils.ParseDatabaseURI(dbURI)
		utils.ErrorPanicPrinter(err, true)

		println("Connecting to production database")
	} else {
		parsedAddress = os.Getenv("DATABASE_URI_DEV")
		println("Connecting to development database")
	}

	db, err := gorm.Open(mysql.Open(parsedAddress), &gorm.Config{})
	utils.ErrorPanicPrinter(err, true)

	err = db.AutoMigrate(&models2.Publisher{}, &models2.Article{}, &models2.User{})
	utils.ErrorPanicPrinter(err, true)

	return db, err
}
