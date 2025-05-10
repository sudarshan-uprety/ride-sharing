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

	response.Success(c, http.StatusCreated, "user registered successfully", res, nil)

}
