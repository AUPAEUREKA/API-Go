package database

import (
	"API-Go/model"
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func Init() {
	db, err := gorm.Open("postgres", "host=127.0.0.1 port=50124 user=apigo dbname=govote password=go-api sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	fmt.Println("Successfully connected!")

	db.AutoMigrate(&model.User{}, &model.Proposal{})
}
