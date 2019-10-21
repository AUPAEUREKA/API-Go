package db

import (
	"API/model"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type PostGresSQL struct {
	db *gorm.DB
}

type Persist interface {
	Connect() (Persist, error)
	SaveUser(model.User) error
	GetUser() (model.User, error)
	GetUsers() ([]model.User, error)
	SaveProposal(model.Proposal) error
	GetProposal() ([]model.Proposal, error)
}

func (p PostGresSQL) Connect() (Persist, error) {
	login := "host=URL docker port=PORT docker user=User Docker dbname=databse Docker password= pass Docker"
	db, err := gorm.Open("postgres", login)
	if err != nil {
		return nil, err
	}
	// Migrate the schema
	db.AutoMigrate(&User{}, &Proposal{})
	p.db = db
	return &p, nil
}

func (p PostGresSQL) SaveUser(u User) error {
	return p.db.Create(&u).Error
}

func (p PostGresSQL) GetUser() (User, error) {
	var us User
	err := p.db.Find(&us).Error
	return us, err
}

func (p PostGresSQL) GetUsers() ([]User, error) {
	var us []User
	err := p.db.Find(&us).Error
	return us, err
}

func (p PostGresSQL) SaveProposal(pr Proposal) error {
	return p.db.Create(&pr).Error
}

func (p PostGresSQL) GetProposal() ([]Proposal, error) {
	var ps []Proposal
	err := p.db.Find(&ps).Error
	return ps, err
}
