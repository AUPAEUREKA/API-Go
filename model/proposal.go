package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Proposal : the proposal struct definition
type Proposal struct {
	gorm.Model
	UUID      string    `json:"uuid"`
	Title     string    `json:"title"`
	Desc      string    `json:"desc"`
	UUIDVote  []User    `gorm:"many2many:ProposalUsers;association_jointable_foreignkey:UUID" json:"uuid_vote"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}
