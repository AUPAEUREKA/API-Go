package main

import (
	"API-Go/database"
	"API-Go/handler"

	"github.com/gin-gonic/gin"
)

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
	database.Init
	r := gin.Default()
	//r.POST("/login", PostLoginHandler)
	r.PUT("/users/:uuid", handler.UpdateUser)
	r.DELETE("/users/:uuid", handler.DeleteUser)
	r.POST("/users", handler.CreateUser)
	//r.PUT("/users/:uuid", PutUserHandler)
	//r.DELETE("/users/:uuid", DeleteUserHandler)
	//r.POST("/votes", PostVoteHandler)
	//r.GET("/votes/:uuid", GetVoteHandler)
	//r.PUT("/votes/:uuid", PutVoteHandler)
	//
	r.Run(":8080")
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
