package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var DB db.Persist

func GetUser(c *gin.Context) {

	us, _ := DB.GetUser()
	c.JSON(http.StatusOK, us)
}
