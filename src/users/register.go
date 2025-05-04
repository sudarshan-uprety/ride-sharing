package users

import (
	"net/http"
	"ride-sharing/initializers"
	"ride-sharing/models"
	userSchemas "ride-sharing/schemas/users"
	"ride-sharing/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(c *gin.Context) {
	// Use the UserRegisterRequest struct instead of AuthInput
	var registerRequest userSchemas.UserRegisterRequest

	if err := c.ShouldBindJSON(&registerRequest); err != nil {
		utils.HandleRequestErrors(c, err)
		return
	}

	// Validation
	if err := registerRequest.Validate(c); err != nil {
		errorMessages := utils.FormatValidatorError(err)
		utils.Error(c.Writer, http.StatusBadRequest, "Validation failed", errorMessages)
		return
	}

	// If we reach here, all validations have passed
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.Error(c.Writer, http.StatusInternalServerError, "Failed to hash password", err.Error())
		return
	}

	// Create user with more fields from the request
	user := models.User{
		Email:    registerRequest.Email,
		FullName: registerRequest.FullName,
		Phone:    registerRequest.Phone,
		Address:  registerRequest.Address,
		Password: string(passwordHash),
	}

	initializers.DB.Create(&user)

	response := userSchemas.UserRegisterResponse{
		ID:        user.ID,
		Email:     user.Email,
		FullName:  user.FullName,
		Phone:     user.Phone,
		Address:   user.Address,
		CreatedAt: user.CreatedAt,
	}

	utils.Success(c.Writer, http.StatusCreated, "User created successfully", response, "")
}
