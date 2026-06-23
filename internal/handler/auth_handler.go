package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/merkurtran/go-im-core/internal/service"
	"github.com/merkurtran/go-im-core/internal/response"
)

type AuthHandler struct {
	userSvc *service.UserService
}

func NewAuthHandler(userSvc *service.UserService) *AuthHandler {
	return &AuthHandler{
		userSvc: userSvc,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 1001, "invalid request")
		return
	}

	resp, err := h.userSvc.Register(c.Request.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserAlreadyExists):
			response.Error(c, http.StatusConflict, 2001, "user already exists")
		default:
			response.Error(c, http.StatusInternalServerError, 5001, "server error")
		}
		return
	}
	response.Success(c, resp)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 1001, "invalid request")
		return
	}

	resp, err := h.userSvc.Login(c.Request.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidUsernameOrPassword):
			response.Error(c, http.StatusUnauthorized, 2002, "invalid username or password")
		default:
			response.Error(c, http.StatusInternalServerError, 5001, "server error")
		}
		return
	}
	response.Success(c, resp)

}
