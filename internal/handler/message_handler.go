package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/merkurtran/go-im-core/internal/response"
	"github.com/merkurtran/go-im-core/internal/service"
)

type MessageHandler struct {
	msgSvc *service.MessageService
}

func NewMessageHandler(msgSvc *service.MessageService) *MessageHandler {
	return &MessageHandler{msgSvc: msgSvc}
}

// SendMessage  POST /messages
func (h *MessageHandler) SendMessage(c *gin.Context) {
	var req service.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 1001, "invalid request")
		return
	}

	// 强制从认证上下文获取发送者，防止伪造
	req.SenderID = c.GetString("user_id")
	if req.SenderID == "" {
		response.Error(c, http.StatusUnauthorized, 1002, "unauthorized")
		return
	}

	resp, err := h.msgSvc.SendMessage(c.Request.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			response.Error(c, http.StatusNotFound, 1004, "user not found")
		default:
			response.Error(c, http.StatusInternalServerError, 5001, "server error")
		}
		return
	}
	response.Success(c, resp)
}

// GetConversation  GET /messages?user_id=xxx&limit=20&offset=0
func (h *MessageHandler) GetConversation(c *gin.Context) {
	currentUserID := c.GetString("user_id")
	targetUserID := c.Query("user_id")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	if targetUserID == "" {
		response.Error(c, http.StatusBadRequest, 1001, "user_id is required")
		return
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	messages, err := h.msgSvc.GetConversation(c.Request.Context(), currentUserID, targetUserID, limit, offset)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 5001, "server error")
		return
	}
	response.Success(c, messages)
}

// MarkRead  PATCH /messages/read
func (h *MessageHandler) MarkRead(c *gin.Context) {
	var req struct {
		OtherUserID string `json:"other_user_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 1001, "invalid request")
		return
	}

	currentUserID := c.GetString("user_id")
	if err := h.msgSvc.MarkConversationAsRead(c.Request.Context(), currentUserID, req.OtherUserID); err != nil {
		response.Error(c, http.StatusInternalServerError, 5001, "server error")
		return
	}
	response.Success(c, nil)
}
