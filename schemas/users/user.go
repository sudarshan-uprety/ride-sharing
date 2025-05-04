package userSchemas

import (
	"errors"
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
	// 1. First validate required fields (email, full_name, phone, etc.)
	if ok, err := govalidator.ValidateStruct(u); !ok {
		// Returns which field is missing/invalid
		utils.Error(c.Writer, utils.ERROR_BAD_REQUEST_CODE, "Bad request", err.Error())
	}
	// 2. Validate password strength (only if passwords match)
	if err := utils.ValidatePassword(u.Password); err != nil {
		utils.Error(c.Writer, utils.ERROR_BAD_REQUEST_CODE, "password validation failed.", err.Error())
	}

	// 3. Check if Password and ConfirmPassword match (early exit if mismatch)
	if u.Password != u.ConfirmPassword {
		return errors.New("password and confirm password must match")
	}

	// 4. Check if email/phone already exists
	if userQueries.EmailExists(u.Email) {
		utils.Error(c.Writer, utils.ERROR_RESOURCE_ALREADY_EXISTS_CODE, utils.ERROR_RESOURCE_ALREADY_EXISTS, nil)

	}
	if userQueries.PhoneExists(u.Phone) {
		utils.Error(c.Writer, utils.ERROR_RESOURCE_ALREADY_EXISTS_CODE, utils.ERROR_RESOURCE_ALREADY_EXISTS, nil)
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
