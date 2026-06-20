package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/merkurtran/go-im-core/internal/service"
	"github.com/merkurtran/go-im-core/pkg/response"
)

type UserHandler struct {
	userSvc *service.UserService
}

// Get /users/me
func (h *UserHandler) GetProfile(c *gin.Context) {
	//
	userID := c.GetString("user_id")
	user, err := h.userSvc.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, 5001, "server error")
		return
	}
	response.Success(c, user)

}

// Patch /users/me
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	var req service.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 1001, "invalid request")
		return
	}

	err := h.userSvc.UpdateUser(c.Request.Context(), userID, &req)
	if err != nil {
		response.Error(c, 5001, "server error")
		return
	}
	response.Success(c, nil)
}
