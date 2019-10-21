package db

import (
	"API/model"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type PostGresSQL struct {
	db *gorm.DB
}

func (p PostGresSQL) init() {
	login := "host=35.187.178.167 port=5432 user=henri dbname=postgres password=admin"
	conn, err := gorm.Open("postgres", login)

	if err != nil {
		panic("failed to connect database")
	}
	
	db.AutoMigrate(
		&model.User{},
	 	&model.Proposal{}
	)
	p.db = db
	return &p, nil
}
