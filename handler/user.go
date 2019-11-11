package handler

import (
	"API-Go/model"
	"fmt"
	"log"
	"net/http"
	"os"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

var dbUser *gorm.DB
var errUser error

// hashPassword : simple password hashing method
func hashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes)
}

// CheckPasswordHash : Compare password with a hash
func CheckPasswordHash(password, hash string) bool {
	errUser := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return errUser == nil
}

// CreateUser : create a user with validations
func CreateUser(c *gin.Context) {
	var user model.User
	errUser := c.BindJSON(&user)
	if errUser != nil {
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
	/*if age.Age(user.DateOfBirth) < 18 {
		log.Println(err)
		c.JSON(http.StatusBadRequest, "You are not adult!")
		return
	}*/
	if !dbUser.Where("email = ?", user.Email).Find(&user).RecordNotFound() {
		c.JSON(http.StatusBadRequest, "User with this email already exist")
		return
	}
	id, _ := uuid.NewV4()
	// 1 = single user; 2 = admin
	user.AccessLevel = 1
	user.UUID = id.String()
	var hash = hashPassword(user.Password)
	user.Password = hash
	dbUser.Create(&user)
	c.JSON(200, &user)
}

// UpdateUser : update a user
func UpdateUser(c *gin.Context) {
	var user model.User
	uuid := c.Params.ByName("uuid")

	if errUser := db.Where("uuid = ?", uuid).First(&user).Error; errUser != nil {
		c.AbortWithStatus(400)
		panic(errUser)
	}
	c.BindJSON(&user)
	var hash = hashPassword(user.Password)
	user.Password = hash
	dbUser.Save(&user)
	c.JSON(200, user)
}

// DeleteUser : delete a user
func DeleteUser(c *gin.Context) {
	var user model.User
	uuid := c.Params.ByName("uuid")

	dbUser.Where("uuid = ?", uuid).Delete(&user)
	c.JSON(200, "deleted")
}

// Login : connect a user and add a JWT token
func Login(c *gin.Context) {
	type login struct {
		Username string `json:"email"`
		Password string `json:"password"`
	}

	loginParams := login{}
	c.ShouldBindJSON(&loginParams)
	type Result struct {
		Email       string
		UUID        string
		Accesslevel string
		Password    string
	}
	var user Result
	if !dbUser.Table("users").Select("email, password, uuid, access_level").Where("email = ?", loginParams.Username).Scan(&user).RecordNotFound() {
		if CheckPasswordHash(loginParams.Password, user.Password) {
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"uuid": user.UUID,
				"acl":  user.Accesslevel,
			})
			tokenStr, errUser := token.SignedString([]byte(os.Getenv("SECRET")))

			if errUser != nil {
				c.JSON(500, errUser)
				return
			}

			c.JSON(200, tokenStr)
			return
		} else {
			c.JSON(http.StatusBadRequest, "wrong password")
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, "wrong email")
		return
	}
}
