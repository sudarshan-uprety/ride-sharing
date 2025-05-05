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
	var registerRequest userSchemas.UserRegisterRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&registerRequest); err != nil {
		utils.HandleRequestErrors(c, err)
		return
	}

	// Custom validation
	if err := registerRequest.Validate(); err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.ErrorResponse(c, utils.NewErrorResponse(
			http.StatusInternalServerError,
			"Failed to hash password",
			err.Error(),
			nil,
		))
		return
	}

	// Create user
	user := models.User{
		Email:    registerRequest.Email,
		FullName: registerRequest.FullName,
		Phone:    registerRequest.Phone,
		Address:  registerRequest.Address,
		Password: string(passwordHash),
	}

	if err := initializers.DB.Create(&user).Error; err != nil {
		utils.ErrorResponse(c, utils.NewErrorResponse(
			http.StatusInternalServerError,
			"Failed to create user",
			err.Error(),
			nil,
		))
		return
	}

	response := userSchemas.UserRegisterResponse{
		ID:        user.ID,
		Email:     user.Email,
		FullName:  user.FullName,
		Phone:     user.Phone,
		Address:   user.Address,
		CreatedAt: user.CreatedAt,
	}

	utils.SuccessResponse(c, http.StatusCreated, "User created successfully", response, "")
}
