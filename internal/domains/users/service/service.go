package service

import (
	"context"
	"errors"
	"ride-sharing/internal/domains/users/dto"
	"ride-sharing/internal/domains/users/models"
	"ride-sharing/internal/domains/users/repository"
	"ride-sharing/internal/pkg/auth"
	customError "ride-sharing/internal/pkg/errors"
	"ride-sharing/internal/pkg/otp"
	"ride-sharing/internal/pkg/password"
	"ride-sharing/internal/pkg/redis"
	"time"
)

type UserService struct {
	repo          repository.UserRepository
	tokenService  *auth.TokenService
	OTPStore      *redis.OTPStore
	userProviders map[auth.UserType]auth.UserProvider
}

func NewUserService(repo repository.UserRepository, tokenService *auth.TokenService, otpStore *redis.OTPStore, userProviders map[auth.UserType]auth.UserProvider) *UserService {
	return &UserService{
		repo:          repo,
		tokenService:  tokenService,
		OTPStore:      otpStore,
		userProviders: userProviders,
	}
}

func (s *UserService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.UserResponse, *customError.AppError) {
	// Check if email exists
	exists, err := s.repo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, customError.NewInternalError(err)
	}
	if exists {
		return nil, customError.NewConflictError("email already exists")
	}

	// Check if phone exists
	exists, err = s.repo.ExistsByPhone(ctx, req.Phone)
	if err != nil {
		return nil, customError.NewInternalError(err)
	}
	if exists {
		return nil, customError.NewConflictError("phone number already exists")
	}

	// Hash password
	hashedPassword, err := password.HashPassword(req.Password)
	if err != nil {
		return nil, customError.NewInternalError(err)
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
		return nil, customError.NewInternalError(err)
	}

	return &dto.UserResponse{
		ID:       user.ID,
		Email:    user.Email,
		FullName: user.FullName,
		Phone:    user.Phone,
	}, nil
}

func (s *UserService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, *customError.AppError) {
	user, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, customError.NewInternalError(err) // Wrap the error
	}
	if user == nil {
		return nil, customError.NewNotFoundError("user not found")
	}

	match, err := password.CheckPassword(req.Password, user.Password)
	if err != nil {
		return nil, customError.NewInternalError(err)
	}
	if !match {
		return nil, customError.NewUnauthorizedError("invalid credentials")
	}

	accessToken, err := s.tokenService.GenerateAccessToken(user.ID.String(), auth.UserTypeUser, user.PasswordChangedAt)
	if err != nil {
		return nil, customError.NewInternalError(err)
	}

	refreshToken, err := s.tokenService.GenerateRefreshToken(user.ID.String(), auth.UserTypeUser, user.PasswordChangedAt)
	if err != nil {
		return nil, customError.NewInternalError(err)
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

func (s *UserService) RefreshToken(ctx context.Context, req dto.RefreshRequest) (*dto.RefreshResponse, *customError.AppError) {
	refreshClaims, err := s.tokenService.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, customError.NewUnauthorizedError("invalid refresh token")
	}

	provider, exists := s.userProviders[refreshClaims.UserType]
	if !exists {
		return nil, customError.NewUnauthorizedError("invalid user type")
	}

	user, err := provider.GetByID(ctx, refreshClaims.UserID, refreshClaims.UserType)
	if err != nil {
		return nil, customError.NewInternalError(err) // Wrap the error
	}
	if user == nil {
		return nil, customError.NewNotFoundError("user not found")
	}

	userData := user.(*models.User)

	tokenPasswordChangedAt := time.Unix(0, refreshClaims.PasswordChangedAt)
	if tokenPasswordChangedAt.Before(*userData.PasswordChangedAt) {
		return nil, customError.NewUnauthorizedError("password changed - please login again")
	}

	accessToken, err := s.tokenService.GenerateAccessToken(userData.ID.String(), auth.UserTypeUser, userData.PasswordChangedAt)
	if err != nil {
		return nil, customError.NewInternalError(err)
	}

	return &dto.RefreshResponse{
		AccessToken: accessToken,
	}, nil

}

func (s *UserService) ChangePassword(ctx context.Context, userID string, req dto.ChangePasswordRequest) (*dto.LoginResponse, *customError.AppError) {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, customError.NewInternalError(err)
	}
	if user == nil {
		return nil, customError.NewNotFoundError("user not found")
	}

	match, err := password.CheckPassword(req.CurrentPassword, user.Password)
	if err != nil {
		return nil, customError.NewInternalError(err)
	}
	if !match {
		return nil, customError.NewVerificationError("incorrect current password")
	}

	hashedPassword, err := password.HashPassword(req.NewPassword)
	if err != nil {
		return nil, customError.NewInternalError(err)
	}

	success, err := s.repo.ChangePassword(ctx, user, hashedPassword)
	if err != nil || !success {
		return nil, customError.NewInternalError(err)
	}

	accessToken, err := s.tokenService.GenerateAccessToken(user.ID.String(), auth.UserTypeUser, user.PasswordChangedAt)
	if err != nil {
		return nil, customError.NewInternalError(err)
	}
	refreshToken, err := s.tokenService.GenerateRefreshToken(user.ID.String(), auth.UserTypeUser, user.PasswordChangedAt)
	if err != nil {
		return nil, customError.NewInternalError(err)
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

func (s *UserService) ForgetPassword(ctx context.Context, req dto.ForgetPasswordRequest) (bool, *customError.AppError) {
	user, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return false, customError.NewInternalError(err)
	}
	if user == nil {
		return false, customError.NewNotFoundError("user not found")
	}

	otp := otp.GenerateOTP()

	if err := s.OTPStore.SetOTP(ctx, user.Email, otp); err != nil {
		var conflictErr *customError.AppError
		if errors.As(err, &conflictErr) && conflictErr.Type == customError.ErrorTypeConflict {
			// Return the conflict error directly
			return false, conflictErr
		}
		return false, customError.NewInternalError(err)
	}
	// remaining: send email functionality
	return true, nil
}

func (s *UserService) VerifyForgetPassword(ctx context.Context, req dto.ForgetPasswordVerifyRequest) (bool, *customError.AppError) {
	user, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return false, customError.NewInternalError(err)
	}
	if user == nil {
		return false, customError.NewNotFoundError("user not found")
	}

	valid, err := s.OTPStore.VerifyAndDeleteOTP(ctx, req.Email, req.Otp)
	if err != nil {
		return false, customError.NewInternalError(err)
	}
	if !valid {
		return false, customError.NewVerificationError("invalid or expired OTP")
	}

	hashedPassword, err := password.HashPassword(req.Password)
	if err != nil {
		return false, customError.NewInternalError(err)
	}

	success, err := s.repo.ChangePassword(ctx, user, hashedPassword)
	if err != nil || !success {
		return false, customError.NewInternalError(err)
	}
	// Remaining: send email
	return true, nil
}
