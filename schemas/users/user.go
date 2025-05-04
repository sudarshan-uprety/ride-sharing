package userSchemas

import (
	"errors"
	"regexp"
	userQueries "ride-sharing/src/users/queries"

	"github.com/asaskevich/govalidator"
)

type UserRegisterRequest struct {
	Email           string `json:"email" validate:"required,email"`
	FullName        string `json:"full_name" validate:"required"`
	Phone           string `json:"phone" validate:"required,min=10"`
	Address         string `json:"address" validate:"required"`
	Password        string `json:"password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required"`
}

func (u *UserRegisterRequest) Validate() error {
	// Validate struct fields based on tags
	if ok, err := govalidator.ValidateStruct(u); !ok {
		return err
	}

	// Custom password rule: min 8, 1 uppercase, 1 digit, 1 special char
	passwordRegex := `^(?=.*[A-Za-z])(?=.*\d)(?=.*[@$!%*#?&])[A-Za-z\d@$!%*#?&]{8,}$`
	matched, _ := regexp.MatchString(passwordRegex, u.Password)
	if !matched {
		return errors.New("password must be at least 8 characters long and include 1 uppercase letter, 1 digit, and 1 special character")
	}

	// Check if Password and ConfirmPassword match
	if u.Password != u.ConfirmPassword {
		return errors.New("Password and Confirm Password must match")
	}

	// Check if email already exists
	if userQueries.EmailExists(u.Email) {
		return errors.New("Email already registered")
	}

	// Check if phone already exists
	if userQueries.PhoneExists(u.Phone) {
		return errors.New("Phone number already registered")
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
