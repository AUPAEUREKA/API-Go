package handler

import (
	"API-Go/model"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

var db *gorm.DB
var err error

//Create a new proposal
func CreateProposal(c *gin.Context) {
	var proposal model.Proposal
	c.BindJSON(&proposal)
	id := uuid.NewV4()
	proposal.UUID = id.String()
	proposal.StartDate = time.Now()
	proposal.EndDate = time.Now().AddDate(0, 3, 0)
	db.Create(&proposal)
	c.JSON(200, &proposal)
}

//Update a proposal
func UpdateProposal(c *gin.Context) {
	var proposal model.Proposal
	uuid := c.Params.ByName("uuid")

	if err := db.Where("uuid = ?", uuid).First(&proposal).Error; err != nil {
		c.AbortWithStatus(400)
		panic(err)
	}
	c.BindJSON(&proposal)
	db.Save(&proposal)
	c.JSON(200, proposal)
}

//User votes a proposal
func Vote(c *gin.Context) {
	var user model.User
	var proposal model.Proposal
	uuidUser := c.Params.ByName("uuid_user")
	uuidProposal := c.Params.ByName("uuid_prop")

	if err := db.Where("uuid = ?", uuidUser).First(&user).Error; err != nil {
		fmt.Println("err user")
		c.AbortWithStatus(400)
		panic(err)
	}
	if err := db.Where("uuid = ?", uuidProposal).First(&proposal).Error; err != nil {
		fmt.Println("err prop")
		c.AbortWithStatus(400)
		panic(err)
	}

	if err := db.Model(&proposal).Association("UUIDVote").Append(&user).Error; err != nil {
		fmt.Println("err query")
		c.AbortWithStatus(400)
	} else {
		c.JSON(200, "Voted")
	}

}

//Get data a proposal
func GetProposal(c *gin.Context) {
	var proposal model.Proposal
	uuid := c.Params.ByName("uuid")

	if err := db.Where("uuid = ?", uuid).First(&proposal).Error; err != nil {
		c.AbortWithStatus(400)
		panic(err)
	}
	if err := db.Preload("UUIDVote").First(&proposal, "uuid = ?", uuid).Error; err != nil {
		c.AbortWithStatus(400)
		panic(err)
	}
	c.JSON(200, &proposal)
}
