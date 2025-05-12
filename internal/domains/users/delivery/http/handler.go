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
