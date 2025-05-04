package userSchemas

import (
	"errors"
	userQueries "ride-sharing/src/users/queries"
	"ride-sharing/utils"

	"github.com/asaskevich/govalidator"
)

type UserRegisterRequest struct {
	Email           string `json:"email" binding:"required,email"`
	FullName        string `json:"full_name" binding:"required"`
	Phone           string `json:"phone" binding:"required,min=10"`
	Address         string `json:"address" binding:"required"`
	Password        string `json:"password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

func (u *UserRegisterRequest) Validate() error {
	// 1. First validate required fields (email, full_name, phone, etc.)
	if ok, err := govalidator.ValidateStruct(u); !ok {
		return err // Returns which field is missing/invalid
	}

	// 2. Check if Password and ConfirmPassword match (early exit if mismatch)
	if u.Password != u.ConfirmPassword {
		return errors.New("password and confirm password must match")
	}

	// 3. Validate password strength (only if passwords match)
	if err := utils.ValidatePassword(u.Password); err != nil {
		return err
	}

	// 4. Check if email/phone already exists
	if userQueries.EmailExists(u.Email) {
		return errors.New("email already registered")
	}
	if userQueries.PhoneExists(u.Phone) {
		return errors.New("phone number already registered")
	}

	return nil
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
