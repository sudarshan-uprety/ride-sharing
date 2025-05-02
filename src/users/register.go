package users

import (
	"net/http"
	"ride-sharing/initializers"
	"ride-sharing/models"
	userSchemas "ride-sharing/schemas/users"
	"ride-sharing/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(c *gin.Context) {

	var authInput userSchemas.AuthInput

	if err := c.ShouldBindJSON(&authInput); err != nil {
		utils.Error(c.Writer, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}

	var userFound models.User
	initializers.DB.Where("email=?", authInput.Email).Find(&userFound)

	if userFound.ID != uuid.Nil {
		utils.Error(c.Writer, http.StatusBadRequest, "Email already used", nil)
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(authInput.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.Error(c.Writer, http.StatusInternalServerError, "Failed to hash password", err.Error())
		return
	}

	user := models.User{
		Email:    authInput.Email,
		Password: string(passwordHash),
	}

	initializers.DB.Create(&user)

	utils.Success(c.Writer, http.StatusCreated, "User created successfully", user, "")

}
