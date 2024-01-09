package database

import (
	"auth-hex/models"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func DB_Init() *gorm.DB {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		"root",
		"",
		"127.0.0.1",
		"3306",
		"hex_str",
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	fmt.Println("Success : Success to connect database.")

	db.AutoMigrate(&models.User{})

	return db
}
