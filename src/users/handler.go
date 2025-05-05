package users

import (
	"net/http"
	"ride-sharing/models"
	userSchemas "ride-sharing/schemas/users"
	queries "ride-sharing/src/users/queries"
	"ride-sharing/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	Repo      queries.UserRepository
	Validator *userSchemas.UserValidator
}

func NewUserHandler(repo queries.UserRepository, validator *userSchemas.UserValidator) *UserHandler {
	return &UserHandler{
		Repo:      repo,
		Validator: validator,
	}
}

func (h *UserHandler) RegisterUser(c *gin.Context) {
	var registerRequest userSchemas.UserRegisterRequest

	if err := c.ShouldBindJSON(&registerRequest); err != nil {
		utils.HandleRequestErrors(c, err)
		return
	}

	if err := h.Validator.Validate(c.Request.Context(), &registerRequest); err != nil {
		utils.ErrorResponse(c, err)
		return
	}

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

	user := models.User{
		Email:    registerRequest.Email,
		FullName: registerRequest.FullName,
		Phone:    registerRequest.Phone,
		Address:  registerRequest.Address,
		Password: string(passwordHash),
	}

	createdUser, err := h.Repo.CreateUser(c.Request.Context(), &user)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	response := userSchemas.UserRegisterResponse{
		ID:        createdUser.ID,
		Email:     createdUser.Email,
		FullName:  createdUser.FullName,
		Phone:     createdUser.Phone,
		Address:   createdUser.Address,
		CreatedAt: createdUser.CreatedAt,
	}

	utils.SuccessResponse(c, http.StatusCreated, "User created successfully", response, "")
}
