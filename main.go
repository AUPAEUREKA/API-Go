package main

import (
	"API-Go/model"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	age "github.com/bearbin/go-age"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

var db *gorm.DB
var err error

func main() {

	//Connect database
	//var db = database.Init
	db, err = gorm.Open("postgres", "host=127.0.0.1 port=5432 user=apigo dbname=govote password=go-api sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	fmt.Println("Successfully connected!")

	//Migrate models
	db.AutoMigrate(&model.User{}, &model.Proposal{})

	//Routers
	r := gin.Default()
	r.POST("/login", Login)
	r.PUT("/users/:uuid", UpdateUser)
	r.DELETE("/users/:uuid", DeleteUser)
	r.POST("/users", CreateUser)
	r.POST("/votes", CreateProposal)
	r.GET("/votes/:uuid", GetProposal)
	r.PUT("/votes/:uuid", UpdateProposal)
	r.POST("/users/:uuid_user/vote/:uuid_prop", Vote)
	r.Run(":8080")
}

//Connect a user and create a JWT token
func Login(c *gin.Context) {
	os.Setenv("API_SECRET", "85ds47")
	type login struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	loginParams := login{}
	c.ShouldBindJSON(&loginParams)
	type Result struct {
		Email       string
		UUID        string
		AccessLevel string
		Password    string
	}
	var user Result
	if !db.Table("users").Select("email, password, uuid, access_level").Where("email = ?", loginParams.Email).Scan(&user).RecordNotFound() {
		if CheckPasswordHash(loginParams.Password, user.Password) {
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"uuid": user.UUID,
				"acl":  user.AccessLevel,
			})
			tokenStr, err := token.SignedString([]byte(os.Getenv("API_SECRET")))

			if err != nil {
				c.JSON(500, err)
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

//Simple password hashing method
func hashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes)
}

//Compare password with a hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
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
