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

	"go.uber.org/zap"
)

type UserService struct {
	repo         repository.UserRepository
	tokenService *auth.TokenService
	OTPStore     *redis.OTPStore
}

func NewUserService(repo repository.UserRepository, tokenService *auth.TokenService, otpStore *redis.OTPStore) *UserService {
	return &UserService{
		repo:         repo,
		tokenService: tokenService,
		OTPStore:     otpStore,
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
	// 1. Validate refresh token structure and signature
	refreshClaims, err := s.tokenService.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, customError.NewUnauthorizedError("invalid refresh token")
	}

	// 2. Verify user type is valid
	provider, exists := s.userProviders[refreshClaims.UserType]
	if !exists {
		return nil, customError.NewUnauthorizedError("invalid user type")
	}

	// 3. Get current user data
	user, appErr := provider.GetByID(ctx, refreshClaims.UserID, refreshClaims.UserType)
	if appErr != nil {
		return nil, appErr
	}
	if user == nil {
		return nil, customError.NewUnauthorizedError("user no longer exists")
	}

	// 4. Check password change timestamp
	tokenPasswordChangedAt := time.Unix(0, refreshClaims.PasswordChangedAt)
	if tokenPasswordChangedAt.Before(*user.PasswordChangedAt) {
		return nil, customError.NewUnauthorizedError("password changed - please login again")
	}

	// 5. Verify refresh token hasn't been invalidated
	isInvalidated, appErr := s.tokenService.ValidateRefreshToken(req.RefreshToken)
	if appErr != nil {
		return nil, appErr
	}
	if isInvalidated {
		return nil, customError.NewUnauthorizedError("refresh token was revoked")
	}

	// 6. Generate new token pair
	tokenPair, appErr := s.tokenService.GenerateAccessToken(user)
	if appErr != nil {
		return nil, customError.NewInternalError(appErr)
	}

	// 7. Invalidate old refresh token (async)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.tokenRepo.InvalidateRefreshToken(ctx, req.RefreshToken); err != nil {
			s.logger.Error("failed to invalidate refresh token",
				zap.String("userID", user.ID),
				zap.Error(err))
		}
	}()

	return &dto.RefreshResponse{
		AccessToken: tokenPair.AccessToken,
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
