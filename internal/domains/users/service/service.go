package service

import (
	"context"
	"ride-sharing/internal/domains/users/dto"
	"ride-sharing/internal/domains/users/models"
	"ride-sharing/internal/domains/users/repository"
	"ride-sharing/internal/pkg/errors"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.UserResponse, *errors.AppError) {
	// Check if email exists
	exists, err := s.repo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}
	if exists {
		return nil, errors.NewConflictError("email already exists")
	}

	// Check if phone exists
	exists, err = s.repo.ExistsByPhone(ctx, req.Phone)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}
	if exists {
		return nil, errors.NewConflictError("phone number already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}

	// Create user
	user := &models.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		FullName: req.FullName,
		Phone:    req.Phone,
		Active:   true,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, errors.NewInternalError(err)
	}

	return &dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FullName:  user.FullName,
		Phone:     user.Phone,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s *UserService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	// Implementation would include JWT token generation
	// Similar pattern to Register
	return nil, nil
}
