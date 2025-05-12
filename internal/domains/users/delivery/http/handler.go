package http

import (
	"net/http"

	"ride-sharing/internal/domains/users/dto"
	"ride-sharing/internal/domains/users/service"
	"ride-sharing/internal/pkg/errors"
	"ride-sharing/internal/pkg/response"
	"ride-sharing/internal/pkg/validation"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// Register godoc
// @Summary      Register a new user
// @Description  Register a new user with email, password, and other details
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request  body  dto.RegisterRequest  true  "User registration data"
// @Success      201      {object}  response.SuccessResponse{dto.UserResponse}  "User registered successfully"
// @Failure      400      {object}  response.ErrorResponse  "Validation error"
// @Failure      409      {object}  response.ErrorResponse  "User already exists"
// @Failure      500      {object}  response.ErrorResponse  "Internal server error"
// @Router       /users/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		details := validation.ProcessValidationError(err)
		response.Error(c, errors.NewValidationError("invalid request body", details))
		return
	}

	user, err := h.service.Register(c.Request.Context(), req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "user registered successfully", user, nil)
}

// Login godoc
// @Summary      Login a user
// @Description  Authenticate user and return access & refresh tokens
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request  body  dto.LoginRequest  true  "User login credentials"
// @Success      200      {object}  response.SuccessResponse{data=dto.LoginResponse}  "Login successful"
// @Failure      400      {object}  response.ErrorResponse  "Validation error"
// @Failure      401      {object}  response.ErrorResponse  "Invalid credentials"
// @Failure      500      {object}  response.ErrorResponse  "Internal server error"
// @Router       /users/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		details := validation.ProcessValidationError(err)
		response.Error(c, errors.NewValidationError("invalid request body", details))
		return
	}

	res, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, "login successful", res, nil)
}

// Refresh godoc
// @Summary      Refresh access token
// @Description  Get new access token using refresh token
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request  body  dto.RefreshRequest  true  "Refresh token"
// @Success      200      {object}  response.SuccessResponse{data=dto.RefreshResponse}  "Token refreshed successfully"
// @Failure      400      {object}  response.ErrorResponse  "Validation error"
// @Failure      500      {object}  response.ErrorResponse  "Internal server error"
// @Router       /users/refresh [post]
func (h *UserHandler) Refresh(c *gin.Context) {
	var req dto.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		details := validation.ProcessValidationError(err)
		response.Error(c, errors.NewValidationError("invalid request body", details))
		return
	}

	res, err := h.service.RefreshToken(c.Request.Context(), req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, "token fetch successfully.", res, nil)
}

// Change Password godoc
// @Summary      Change user password
// @Description  Change password for authenticated user
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body  dto.ChangePasswordRequest  true  "Change password data"
// @Success      200      {object}  response.SuccessResponse{data=dto.LoginResponse}  "Password changed successfully"
// @Failure      400      {object}  response.ErrorResponse  "Validation error"
// @Failure      401      {object}  response.ErrorResponse  "Unauthorized"
// @Failure      500      {object}  response.ErrorResponse  "Internal server error"
// @Router       /users/change-password [post]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	var req dto.ChangePasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		details := validation.ProcessValidationError(err)
		response.Error(c, errors.NewValidationError("invalid request body", details))
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, errors.NewUnauthorizedError("user ID not found in context"))
		return
	}

	res, err := h.service.ChangePassword(c.Request.Context(), userID.(string), req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, "password changed successfully", res, nil)
}

// Forget Password godoc
// @Summary      Forget password
// @Description  Forget password
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request  body  dto.ForgetPasswordRequest  true  "Forget password data"
// @Success      202      {object}  response.SuccessResponse{data=bool}  "OTP sent to registered mail"
// @Failure      400      {object}  response.ErrorResponse  "Validation error"
// @Failure      500      {object}  response.ErrorResponse  "Internal server error"
// @Router       /users/forget-password [post]
func (h *UserHandler) ForgetPassword(c *gin.Context) {
	var req dto.ForgetPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		details := validation.ProcessValidationError(err)
		response.Error(c, errors.NewValidationError("invalid request body", details))
		return
	}
	res, err := h.service.ForgetPassword(c.Request.Context(), req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusAccepted, "OTP sent to user email", res, nil)
}

// Verify Forget Password godoc
// @Summary      Forget password
// @Description  Forget password
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request  body  dto.ForgetPasswordVerifyRequest  true  "Verify forget password data"
// @Success      200      {object}  response.SuccessResponse{data=bool}  "Password reset successfully."
// @Failure      400      {object}  response.ErrorResponse  "Validation error"
// @Failure      500      {object}  response.ErrorResponse  "Internal server error"
// @Router       /users/verify-reset [post]
func (h *UserHandler) VerifyForgetPassword(c *gin.Context) {
	var req dto.ForgetPasswordVerifyRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		details := validation.ProcessValidationError(err)
		response.Error(c, errors.NewValidationError("invalid request body", details))
		return
	}
	res, err := h.service.VerifyForgetPassword(c.Request.Context(), req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusAccepted, "password reset successfully", res, nil)
}
