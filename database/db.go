package db

import (
	"API-GO/model"
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB

func Init() {
	db, err := gorm.Open("postgres", "host=127.0.0.1 port=5432 user=apigo dbname=govote password=go-api sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	fmt.Println("Successfully connected!")

	db.AutoMigrate(&model.User{}, &model.Proposal{})
}
