package users

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/istomin10593/bookstore_users-api/domain/users"

	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
	var user users.User
	bytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return
	}
	if err := json.Unmarshal(bytes, &user); err != nil {
		return
	}
	c.String(http.StatusNotImplemented, "implement me!")
}

func GetUser(c *gin.Context) {
	c.String(http.StatusNotImplemented, "implement me!")
}

// func SearchUser(c *gin.Context) {
// 	c.String(http.StatusNotImplemented, "implement me!")
// }
