package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/merkurtran/go-im-core/internal/response"
	"github.com/merkurtran/go-im-core/internal/service"
)

type AuthHandler struct {
	userSvc *service.UserService
}

func NewAuthHandler(userSvc *service.UserService) *AuthHandler {
	return &AuthHandler{
		userSvc: userSvc,
	}
}

// Register 注册新用户
//
//	@Summary 注册新用户
//	@Description 用用户名和密码注册新用户，返回 user_id 和 token
//	@Tags 认证
//	@Accept json
//	@Produce json
//	@Param request body service.RegisterRequest true "注册信息"
//	@Success 200 {object} response.Response{data=service.RegisterResponse} "注册成功"
//	@Failure 400 {object} response.Response "请求参数错误"
//	@Failure 409 {object} response.Response "用户已存在"
//	@Failure 500 {object} response.Response "服务器错误"
//	@Router /auth/register [post]
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

// Login 登录
//
//	@Summary 用户登录
//	@Description 用户名和密码登录，返回token
//	@Tags 认证
//	@Accept json
//	@Produce json
//	@Param request body service.LoginRequest true "登录信息"
//	@Success 200 {object} response.Response{data=service.LoginResponse} "登录成功"
//	@Failure 400 {object} response.Response "请求参数错误"
//	@Failure 401 {object} response.Response "用户名或密码错误"
//	@Failure 500 {object} response.Response "服务器错误"
//	@Router /auth/login [post]
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
