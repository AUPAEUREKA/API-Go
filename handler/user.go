package handler

import (
	"API-Go/model"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

var db *gorm.DB
var err error

// hashPassword : simple password hashing method
func hashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes)
}

// CheckPasswordHash : Compare password with a hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// CreateUser : create a user with validations
func CreateUser(c *gin.Context) {
	var user model.User
	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	if len(user.Password) == 0 {
		fmt.Println("err2")
		log.Println(err)
		c.JSON(http.StatusBadRequest, "No given password")
		return
	}
	if user.DateOfBirth < 18 {
		log.Println(err)
		c.JSON(http.StatusBadRequest, "You are not adult!")
		return
	}
	if !db.Where("email = ?", user.Email).Find(&user).RecordNotFound() {
		c.JSON(http.StatusBadRequest, "User with this email already exist")
		return
	}
	id, _ := uuid.NewV4()
	user.AccessLevel = 1
	user.UUID = id.String()
	var hash = hashPassword(user.Password)
	user.Password = hash
	db.Create(&user)
	c.JSON(200, &user)
}

// UpdateUser : update a user
func UpdateUser(c *gin.Context) {
	var user model.User
	uuid := c.Params.ByName("uuid")

	if err := db.Where("uuid = ?", uuid).First(&user).Error; err != nil {
		c.AbortWithStatus(400)
		panic(err)
	}
	c.BindJSON(&user)

	db.Save(&user)
	c.JSON(200, user)
}

// DeleteUser : delete a user
func DeleteUser(c *gin.Context) {
	var user model.User
	uuid := c.Params.ByName("uuid")
	if err := db.Where("uuid = ?", uuid).Delete(&user); err != nil {
		c.AbortWithStatus(400)
		panic(err)
	}
	c.JSON(200, user)
}

// Login : connect a user
func Login(c *gin.Context) {
	var user model.User
	c.BindJSON(&user)
	var pass = user.Password
	fmt.Println(hashPassword(pass))
	if err := db.Where("email = ?", user.Email).First(&user).Error; err != nil {
		//if CheckPasswordHash(hashedPassword, user.Password)
		c.JSON(404, "Wrong email")
		fmt.Println(err)
	} else {
		db.Table("users").Select("password").Where("email = ?", user.Email).Scan(&user)
		verify := CheckPasswordHash(pass, user.Password)
		if verify {
			c.JSON(200, user)
		} else {
			c.JSON(404, "Wrong password")
		}

	}
}
