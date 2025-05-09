package provider

import (
	"context"
	"fmt"
	"ride-sharing/internal/domains/users/repository"
	"ride-sharing/internal/pkg/auth"
)

// UserProvider now works with multiple user types
type UserProvider struct {
	userRepo repository.UserRepository
	// adminRepo repository.AdminRepository
	// riderRepo repository.RiderRepository
}

// NewUserProvider initializes the UserProvider with necessary repositories
func NewUserProvider(userRepo repository.UserRepository) auth.UserProvider {
	return &UserProvider{
		userRepo: userRepo,
		// adminRepo: adminRepo,
		// riderRepo: riderRepo,
	}
}

// GetByID fetches a user based on the UserType (User, Admin, Rider)
func (p *UserProvider) GetByID(ctx context.Context, id string, userType auth.UserType) (interface{}, error) {
	switch userType {
	// case auth.UserTypeAdmin:
	// 	// If the user is an admin, use the admin repository
	// 	return p.adminRepo.GetByID(ctx, id)
	// case auth.UserTypeRider:
	// 	// If the user is a rider, use the rider repository
	// 	return p.riderRepo.GetByID(ctx, id)
	case auth.UserTypeUser:
		// If the user is a regular user, use the user repository
		return p.userRepo.GetByID(ctx, id)
	default:
		return nil, fmt.Errorf("invalid user type: %s", userType)
	}
}
