package model

import (
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// Proposal : the survey struct definition
type Proposal struct {
	UUID      string    `json:"uuid"`
	Title     string    `json:"title"`
	Desc      string    `json:"desc"`
	UUIDVote  []User    `json:"uuid_vote"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName : Gorm
func (p *Proposal) TableName() string {
	return "Proposals"
}

// BeforeCreate : Gorm hook
func (p *Proposal) BeforeCreate(scope *gorm.Scope) {
	id, _ := uuid.NewV4()
	p.UUID = id.String()
	scope.SetColumn("UUID", p.UUID)
	return
}
