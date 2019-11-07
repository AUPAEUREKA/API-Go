package main

import (
	"API-Go/model"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

var db *gorm.DB
var err error

// User is the representation of a client.
/*type User struct {
	gorm.Model
	UUID        string    `json:"uuid"`
	AccessLevel int       `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	DateOfBirth time.Time `json:"birth_date"`
}

// Vote is the representation of a vote
type Vote struct {
	gorm.Model
	UUID      string    `json:"uuid"`
	Title     string    `json:"title"`
	Desc      string    `json:"desc"`
	UUIDVote  []User    `json:"uuid_vote"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

var listUser = map[string]User{}*/

func main() {

	db, err = gorm.Open("postgres", "host=127.0.0.1 port=5432 user=apigo dbname=govote password=go-api sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	fmt.Println("Successfully connected!")

	db.AutoMigrate(&model.User{} /*, &model.Proposal{}*/)

	//var db = database.Init
	r := gin.Default()
	r.POST("/login", Login)
	r.PUT("/users/:uuid", UpdateUser)
	r.DELETE("/users/:uuid", DeleteUser)
	r.POST("/users", CreateUser)
	//r.PUT("/users/:uuid", PutUserHandler)
	//r.DELETE("/users/:uuid", DeleteUserHandler)
	//r.POST("/votes", PostVoteHandler)
	//r.GET("/votes/:uuid", GetVoteHandler)
	//r.PUT("/votes/:uuid", PutVoteHandler)
	//
	r.Run(":8080")
}

func Login(c *gin.Context) {
	type login struct {
		Username string `json:"email"`
		Password string `json:"password"`
	}

	loginParams := login{}
	//var user model.User
	c.ShouldBindJSON(&loginParams)
	//c.BindJSON(&user)
	type Result struct {
		Email       string
		UUID        string
		Accesslevel string
		Password    string
	}
	var user Result
	if !db.Table("users").Select("email, password, uuid, access_level").Where("email = ?", loginParams.Username).Scan(&user).RecordNotFound() {
		if CheckPasswordHash(loginParams.Password, user.Password) {
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"uuid": user.UUID,
				"acl":  user.Accesslevel,
			})
			tokenStr, err := token.SignedString([]byte(os.Getenv("SECRET")))

			if err != nil {
				c.JSON(500, err)
				return
			}

			c.JSON(200, tokenStr)
			return
		}
	}
	/*
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

		}*/
}

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
	// 1 = single user; 2 = admin
	user.AccessLevel = 1
	user.UUID = id.String()
	var hash = hashPassword(user.Password)
	user.Password = hash
	db.Create(&user)
	c.JSON(200, &user)
}

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
	c.JSON(200, user)
}

func DeleteUser(c *gin.Context) {
	var user model.User
	uuid := c.Params.ByName("uuid")

	db.Where("uuid = ?", uuid).Delete(&user)
	c.JSON(200, "deleted")
}

// GetUserHandler is retriving user from the given uuid param.
/*func GetUserHandler(ctx *gin.Context) {
	if u, ok := listUser[ctx.Param("uuid")]; ok {
		ctx.JSON(http.StatusOK, u)
		return
	}
	ctx.JSON(http.StatusNotFound, nil)
}

// GetAllUserHandler is retriving all users from the database.
func GetAllUserHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, listUser)
}

// PostUserHandler is creating a new user into the database.
func PostUserHandler(ctx *gin.Context) {
	var u User
	if err := ctx.BindJSON(&u); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	u.UUID = uuid.NewV4().String()
	listUser[u.UUID] = u
	ctx.JSON(http.StatusOK, u)
}
*/
