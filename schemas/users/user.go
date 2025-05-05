package userSchemas

import (
	"context"
	"net/http"
	userQueries "ride-sharing/src/users/queries"
	"ride-sharing/utils"
	"time"

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

type UserValidator struct {
	UserRepo userQueries.UserRepository
}

func NewUserValidator(repo userQueries.UserRepository) *UserValidator {
	return &UserValidator{UserRepo: repo}
}

func (v *UserValidator) Validate(ctx context.Context, u *UserRegisterRequest) error {
	// Basic field validation
	if err := utils.ValidatePassword(u.Password); err != nil {
		return utils.NewErrorResponse(
			http.StatusBadRequest,
			"Password validation failed",
			err.Error(),
			nil,
		)
	}

	if u.Password != u.ConfirmPassword {
		return utils.NewErrorResponse(
			http.StatusBadRequest,
			"Password mismatch",
			"Password and confirm password must match",
			nil,
		)
	}

	// Database checks
	emailExists, err := v.UserRepo.EmailExists(ctx, u.Email)
	if err != nil {
		return utils.NewErrorResponse(
			http.StatusInternalServerError,
			"Database error",
			err.Error(),
			nil,
		)
	}
	if emailExists {
		return utils.NewErrorResponse(
			http.StatusConflict,
			"Email already exists",
			"",
			nil,
		)
	}

	phoneExists, err := v.UserRepo.PhoneExists(ctx, u.Phone)
	if err != nil {
		return utils.NewErrorResponse(
			http.StatusInternalServerError,
			"Database error",
			err.Error(),
			nil,
		)
	}
	if phoneExists {
		return utils.NewErrorResponse(
			http.StatusConflict,
			"Phone number already exists",
			"",
			nil,
		)
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
