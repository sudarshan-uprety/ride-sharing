package http

import (
	"net/http"

	"ride-sharing/internal/domains/users/dto"
	"ride-sharing/internal/domains/users/service"
	"ride-sharing/internal/pkg/errors"
	"ride-sharing/internal/pkg/response"

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
		response.Error(c, errors.NewValidationError("invalid request body", err.Error()))
		return
	}

	user, err := h.service.Register(c.Request.Context(), req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "user registered successfully", user, nil)
}

// func (h *UserHandler) Login(c *gin.Context) {
// 	var req dto.LoginRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		response.Error(c, errors.NewValidationError("invalid request body", err.Error()))
// 		return
// 	}

// 	res, err := h.service.Login(c.Request.Context(), req)
// 	if err != nil {
// 		response.Error(c, err)
// 		return
// 	}

// 	response.Success(c, http.StatusOK, "login successful", res, nil)
// }
