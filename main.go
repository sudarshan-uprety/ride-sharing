package main

import (
	"ride-sharing/initializers"
	"ride-sharing/src/users"

	// "ride-sharing/src/users"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvs()
	initializers.ConnectDB()

}

func main() {
	router := gin.Default()

	router.POST("/auth/signup", users.CreateUser)
	router.POST("/auth/login", users.Login)
	// router.GET("/user/profile", middlewares.CheckAuth, users.GetUserProfile)
	router.Run()
}
