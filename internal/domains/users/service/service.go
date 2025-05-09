package service

import (
	"context"
	"ride-sharing/internal/domains/users/dto"
	"ride-sharing/internal/domains/users/models"
	"ride-sharing/internal/domains/users/repository"
	"ride-sharing/internal/pkg/auth"
	"ride-sharing/internal/pkg/errors"
	"ride-sharing/internal/pkg/password"
	"time"
)

type UserService struct {
	repo         repository.UserRepository
	tokenService *auth.TokenService
}

func NewUserService(repo repository.UserRepository, tokenService *auth.TokenService) *UserService {
	return &UserService{
		repo:         repo,
		tokenService: tokenService,
	}
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
	hashedPassword, err := password.HashPassword(req.Password)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}

	// Create user
	current_time := time.Now()
	user := &models.User{
		Email:             req.Email,
		Password:          string(hashedPassword),
		FullName:          req.FullName,
		Phone:             req.Phone,
		Active:            false,
		PasswordChangedAt: &current_time,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, errors.NewInternalError(err)
	}

	return &dto.UserResponse{
		ID:       user.ID,
		Email:    user.Email,
		FullName: user.FullName,
		Phone:    user.Phone,
	}, nil
}

func (s *UserService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, *errors.AppError) {
	user, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.NewInternalError(err) // Wrap the error
	}
	if user == nil {
		return nil, errors.NewNotFoundError("user not found")
	}

	match, err := password.CheckPassword(req.Password, user.Password)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}
	if !match {
		return nil, errors.NewUnauthorizedError("invalid credentials")
	}

	accessToken, err := s.tokenService.GenerateAccessToken(user.ID.String(), auth.UserTypeUser, user.PasswordChangedAt)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}

	refreshToken, err := s.tokenService.GenerateRefreshToken(user.ID.String(), auth.UserTypeUser, user.PasswordChangedAt)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: dto.UserResponse{
			ID:       user.ID,
			Email:    user.Email,
			FullName: user.FullName,
			Phone:    user.Phone,
		},
	}, nil
}

func (s *UserService) ChangePassword(ctx context.Context, userID string, req dto.ChangePasswordRequest) (*dto.LoginResponse, *errors.AppError) {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}
	if user == nil {
		return nil, errors.NewNotFoundError("user not found")
	}

	match, err := password.CheckPassword(req.CurrentPassword, user.Password)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}
	if !match {
		return nil, errors.NewVerificationError("incorrect current password")
	}

	hashedPassword, err := password.HashPassword(req.NewPassword)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}

	success, err := s.repo.ChangePassword(ctx, user, hashedPassword)
	if err != nil || !success {
		return nil, errors.NewInternalError(err)
	}

	accessToken, err := s.tokenService.GenerateAccessToken(user.ID.String(), auth.UserTypeUser, user.PasswordChangedAt)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}
	refreshToken, err := s.tokenService.GenerateRefreshToken(user.ID.String(), auth.UserTypeUser, user.PasswordChangedAt)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: dto.UserResponse{
			ID:       user.ID,
			Email:    user.Email,
			FullName: user.FullName,
			Phone:    user.Phone,
		},
	}, nil
}
