package userSchemas

import (
	"net/http"
	userQueries "ride-sharing/src/users/queries"
	"ride-sharing/utils"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserRegisterRequest struct {
	Email           string `json:"email" binding:"required,email"`
	FullName        string `json:"full_name" binding:"required"`
	Phone           string `json:"phone" binding:"required,min=10"`
	Address         string `json:"address" binding:"required"`
	Password        string `json:"password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

func (u *UserRegisterRequest) Validate(c *gin.Context) error {
	if ok, _ := govalidator.ValidateStruct(u); !ok {
		return utils.Error("Validation error", "Invalid or missing fields", "", http.StatusBadRequest)
	}

	if err := utils.ValidatePassword(u.Password); err != nil {
		return utils.Error("Password validation failed", err.Error(), "", http.StatusBadRequest)
	}

	if u.Password != u.ConfirmPassword {
		return utils.Error("Password mismatch", "Password and confirm password must match", "", http.StatusBadRequest)
	}

	if userQueries.EmailExists(u.Email) {
		return utils.Error("Email already exists", nil, "", http.StatusConflict)
	}

	if userQueries.PhoneExists(u.Phone) {
		return utils.Error("Phone number already exists", nil, "", http.StatusConflict)
	}

	return nil
}

type UserRegisterResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	Phone     string    `json:"phone"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
}

type LoginRequest struct {
	Email    string `json:"email" valid:"required,email"`
	Password string `json:"password" valid:"required,stringlength(8|50)"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AuthInput struct {
	Email    string `json:"email" valid:"required,email"`
	Password string `json:"password" valid:"required,stringlength(6|100)"`
}

// type RegisterRequest struct {
// }

// type RegisterResponse struct {
// }
