package di

import (
	"enchantech-codex/src/models"
	"enchantech-codex/src/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func SetupDatabase() (*gorm.DB, error) {
	//address := "admin:1234@tcp(localhost:11306)/db"
	//connectionString := "mysql://doadmin:AVNS_z0OORLs1Abp8Vbun_0u@enchantech-cluster-do-user-12948347-0.c.db.ondigitalocean.com:25060/enchantech-db?ssl-mode=REQUIRED"

	address := "doadmin:AVNS_z0OORLs1Abp8Vbun_0u@tcp(enchantech-cluster-do-user-12948347-0.c.db.ondigitalocean.com:25060)/enchantech-db?tls=skip-verify"
	db, err := gorm.Open(mysql.Open(address), &gorm.Config{})
	utils.ErrorPanicPrinter(err, true)

	err = db.AutoMigrate(&models.Publisher{}, &models.Article{})
	utils.ErrorPanicPrinter(err, true)

	return db, err
}
