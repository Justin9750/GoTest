package sqlserver

import (
	"fmt"
	"forms/models"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var Database *gorm.DB

func Init() {
	dsn := "sqlserver://s21+:diehards21+@172.24.25.62?database=goTest"
	var err error
	Database,
		err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{})

	Database.Assign(&models.User{})
	// Database.AutoMigrate(&models.User{})
	if err != nil {
		fmt.Println(err.Error())
	}
}
