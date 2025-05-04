package userQueries

import (
	"ride-sharing/initializers"
	"ride-sharing/models"

	"gorm.io/gorm"
)

func EmailExists(email string) bool {
	var userFound models.User
	result := initializers.DB.Where("email = ?", email).First(&userFound)

	// Return true if email is found, false if not found, and handle errors
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return false
		}
		// Optionally, handle other errors here, if needed
		return false
	}
	return true // Email exists
}

func PhoneExists(email string) bool {
	var phoneFound models.User
	result := initializers.DB.Where("phone = ?", email).First(&phoneFound)

	// Return true if phone is found, false if not found, and handle errors
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return false
		}
		// Optionally, handle other errors here, if needed
		return false
	}
	return true // Phone exists
}
