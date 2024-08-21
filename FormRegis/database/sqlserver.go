package sqlserver

import (
	"fmt"
	"forms/models"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var Database *gorm.DB

func Init() {
	dsn := "sqlserver://<username>:<password>@<servername>?database=<dbname>"
	var err error
	Database,
		err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{})

	Database.Assign(&models.User{})
	// Database.AutoMigrate(&models.User{})
	if err != nil {
		fmt.Println(err.Error())
	}
}
