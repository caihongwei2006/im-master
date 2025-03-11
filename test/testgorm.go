package main

import (
	"fmt"
	"im-master/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:cai2006@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	fmt.Println("DSN:", dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}
	db.AutoMigrate(&models.Message{})

}
