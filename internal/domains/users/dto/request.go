// request.go
package dto

import (
	"github.com/google/uuid"
)

type RegisterRequest struct {
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,strongpassword"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
	FullName        string `json:"full_name" binding:"required"`
	Phone           string `json:"phone" binding:"required,e164"`
	Address         string `json:"address" binding:"required"`
}

type UserResponse struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	FullName string    `json:"full_name"`
	Phone    string    `json:"phone"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,strongpassword"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RefreshResponse struct {
	AccessToken string `json:"access_token"`
}

type LoginResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         UserResponse `json:"user"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required,strongpassword"`
	NewPassword     string `json:"new_password" binding:"required,strongpassword"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=NewPassword"`
}

type ForgetPasswordRequest struct {
	Email string `json:"email" binding:"required"`
}

type ForgetPasswordVerifyRequest struct {
	Email           string `json:"email" binding:"required"`
	Otp             string `json:"otp" binding:"required,otpvalidation"`
	Password        string `json:"password" binding:"required,strongpassword"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
}

type VerifyEmailRequest struct {
	Email string `json:"email" binding:"required"`
	Otp   string `json:"otp" binding:"required,otpvalidation"`
}
