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

//Simple password hashing method
func hashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes)
}

//Compare password with a hash
func CheckPasswordHash(password, hash string) bool {
	errUser := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return errUser == nil
}

//Create a user with validations
func CreateUser(c *gin.Context) {
	type result struct {
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		Email       string `json:"email"`
		Password    string `json:"password"`
		DateOfBirth string `json:"birth_date"`
	}
	UserParams := result{}

	err := c.ShouldBindJSON(&UserParams)
	layout := "2006-01-02"
	str := UserParams.DateOfBirth
	t, er := time.Parse(layout, str)

	if er != nil {
		fmt.Println(er)
	}

	var user model.User
	//err := c.BindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	if len(UserParams.Password) == 0 {
		fmt.Println("err2")
		log.Println(err)
		c.JSON(http.StatusBadRequest, "No given password")
		return
	}
	if age.Age(t) < 18 {
		log.Println(err)
		c.JSON(http.StatusBadRequest, "You are not adult!")
		return
	}
	if !db.Where("email = ?", UserParams.Email).Find(&user).RecordNotFound() {
		c.JSON(http.StatusBadRequest, "User with this email already exist")
		return
	}
	id := uuid.NewV4()
	// 1 = single user; 2 = admin
	user.AccessLevel = 1
	user.UUID = id.String()
	var hash = hashPassword(UserParams.Password)
	user.Password = hash
	user.FirstName = UserParams.FirstName
	user.LastName = UserParams.LastName
	user.Email = UserParams.Email
	user.DateOfBirth = t
	db.Create(&user)
	user.Password = ""
	c.JSON(200, &user)
}

//Update a user
func UpdateUser(c *gin.Context) {
	var user model.User
	uuid := c.Params.ByName("uuid")

	if err := db.Where("uuid = ?", uuid).First(&user).Error; err != nil {
		c.AbortWithStatus(400)
		panic(err)
	}
	c.BindJSON(&user)
	var hash = hashPassword(user.Password)
	user.Password = hash
	db.Save(&user)
	user.Password = ""
	c.JSON(200, user)
}

//Delete a user
func DeleteUser(c *gin.Context) {
	var user model.User
	uuid := c.Params.ByName("uuid")

	db.Where("uuid = ?", uuid).Delete(&user)
	c.JSON(200, "This user deletes!")
}

//Connect a user and add a JWT token
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
