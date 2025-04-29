package main

import (
	"ride-sharing/initializers"
	"ride-sharing/models"
)

func init() {
	initializers.LoadEnvs()
	initializers.ConnectDB()

}

func main() {

	initializers.DB.AutoMigrate(&models.User{})
}
