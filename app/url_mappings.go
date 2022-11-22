package app

import (
	"github.com/istomin10593/bookstore_users-api/controllers"
)

func mapUrls() {
	router.GET("/ping", controllers.Ping)

	router.GET("/users/:user_id", controllers.GetUser)
	// router.GET("/users/search", controllers.SearchUser)
	router.POST("/users", controllers.CreateUser)
}
