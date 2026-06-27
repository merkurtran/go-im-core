package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/merkurtran/go-im-core/internal/domain/user"
	"github.com/merkurtran/go-im-core/internal/service"
	"github.com/merkurtran/go-im-core/internal/response"
)

// swaggerDocTypes 用于 swaggo 类型解析，无实际运行效果
var _ = (*user.User)(nil)

type UserHandler struct {
	userSvc *service.UserService
}

func NewUserHandler(userSvc *service.UserService) *UserHandler {
	return &UserHandler{userSvc: userSvc}
}

// GetProfile 获取当前用户信息
//
//	@Summary 获取当前用户信息
//	@Description 获取当前登录用户的个人资料
//	@Tags 用户
//	@Produce json
//	@Success 200 {object} response.Response{data=user.User} "成功"
//	@Failure 404 {object} response.Response "用户不存在"
//	@Failure 500 {object} response.Response "服务器错误"
//	@Security Bearer
//	@Router /users/me [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	user, err := h.userSvc.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 5001, "server error")
		return
	}
	if user == nil {
		response.Error(c, http.StatusNotFound, 1004, "user not found")
		return
	}
	response.Success(c, user)
}

// UpdateProfile 更新当前用户信息
//
//	@Summary 更新个人信息
//	@Description 更新当前用户的昵称或头像
//	@Tags 用户
//	@Accept json
//	@Produce json
//	@Param request body service.UpdateUserRequest true "更新信息"
//	@Success 200 {object} response.Response "更新成功"
//	@Failure 400 {object} response.Response "请求参数错误"
//	@Failure 500 {object} response.Response "服务器错误"
//	@Security Bearer
//	@Router /users/me [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	var req service.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 1001, "invalid request")
		return
	}

	err := h.userSvc.UpdateUser(c.Request.Context(), userID, &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 5001, "server error")
		return
	}
	response.Success(c, nil)
}
