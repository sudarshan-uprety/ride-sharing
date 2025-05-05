package users

// import (
// 	"net/http"
// 	"ride-sharing/models"
// 	userSchemas "ride-sharing/schemas/users"
// 	userQueries "ride-sharing/src/users/queries"
// 	"ride-sharing/utils"

// 	"github.com/gin-gonic/gin"
// 	"golang.org/x/crypto/bcrypt"
// 	"gorm.io/gorm"
// )

// type UserHandler struct {
// 	Validator *userSchemas.UserValidator
// 	Repo      userQueries.UserRepository
// }

// func NewUserHandler(db *gorm.DB) *UserHandler {
// 	repo := userQueries.NewUserRepository(db)
// 	validator := userSchemas.NewUserValidator(repo)
// 	return &UserHandler{
// 		Validator: validator,
// 		Repo:      repo,
// 	}
// }

// func (h *UserHandler) RegisterUser(c *gin.Context) {
// 	var registerRequest userSchemas.UserRegisterRequest

// 	if err := c.ShouldBindJSON(&registerRequest); err != nil {
// 		utils.HandleRequestErrors(c, err)
// 		return
// 	}

// 	// Validate with context
// 	if err := h.Validator.Validate(c.Request.Context(), &registerRequest); err != nil {
// 		utils.ErrorResponse(c, err)
// 		return
// 	}

// 	passwordHash, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), bcrypt.DefaultCost)
// 	if err != nil {
// 		utils.ErrorResponse(c, utils.NewErrorResponse(
// 			http.StatusInternalServerError,
// 			"Failed to hash password",
// 			err.Error(),
// 			nil,
// 		))
// 		return
// 	}

// 	user := models.User{
// 		Email:    registerRequest.Email,
// 		FullName: registerRequest.FullName,
// 		Phone:    registerRequest.Phone,
// 		Address:  registerRequest.Address,
// 		Password: string(passwordHash),
// 	}

// 	if _, err := h.Repo.CreateUser(c.Request.Context(), &user); err != nil {
// 		utils.ErrorResponse(c, err)
// 		return
// 	}

// 	response := userSchemas.UserRegisterResponse{
// 		ID:        user.ID,
// 		Email:     user.Email,
// 		FullName:  user.FullName,
// 		Phone:     user.Phone,
// 		Address:   user.Address,
// 		CreatedAt: user.CreatedAt,
// 	}

// 	utils.SuccessResponse(c, http.StatusCreated, "User created successfully", response, "")
// }
