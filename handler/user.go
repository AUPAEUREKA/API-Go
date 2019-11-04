package handler

import (
	"API-Go/model"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB
var err error

func CreateUser(c *gin.Context) {
	var user model.User
	c.BindJSON(&user)
	db.Create(&user)
	c.JSON(200, user)
}

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

func DeleteUser(c *gin.Context) {
	var user model.User
	uuid := c.Params.ByName("uuid")
	if err := db.Where("uuid = ?", uuid).Delete(&user); err != nil {
		c.AbortWithStatus(400)
		panic(err)
	}
	c.JSON(200, user)
}

/*func NewUser(db *gorm.DB) *UserHandler {
return &UserHandler{db: db}
/*db, err := gorm.Open("postgres", "host=127.0.0.1 port=5432 user=apigo dbname=govote password=go-api sslmode=disable")
if err != nil {
	panic(err)
}
defer db.Close()

vars := mux.Vars(r)
name := vars["name"]
email := vars["email"]
db.Find(&users)
fmt.Println("{}", users)

json.NewEncoder(w).Encode(users)*/
/*}

func (u *UserHandler) Save(user *model.User) UserHandler {
	err := u.db.Save(user).Error

	if err != nil {
		panic(err)
	}

	return UserHandler{User: user}
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Delete User Endpoint Hit")
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Update User Endpoint Hit")
}

func GetUser(c *gin.Context) {

	/*us, _ := DB.GetUser()
	c.JSON(http.StatusOK, us)*/
//}
