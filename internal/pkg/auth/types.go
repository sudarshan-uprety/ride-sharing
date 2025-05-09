package auth

import "context"

type UserType string

const (
	UserTypeAdmin UserType = "admin"
	UserTypeUser  UserType = "user"
	UserTypeRider UserType = "rider"
)

func (u UserType) IsValid() bool {
	switch u {
	case UserTypeAdmin, UserTypeUser, UserTypeRider:
		return true
	}
	return false
}

type UserProvider interface {
	GetByID(ctx context.Context, id string, userType UserType) (interface{}, error)
}
