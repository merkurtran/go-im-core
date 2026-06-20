package handler

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/merkurtran/go-im-core/internal/service"
	"github.com/merkurtran/go-im-core/pkg/response"
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
		response.Error(c, 1001, "invalid request")
		return
	}

	req.SenderID = c.GetString("user_id")

	resp, err := h.msgSvc.SendMessage(c.Request.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			response.Error(c, 1004, "user not found")
		default:
			response.Error(c, 5001, "server error")
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
		response.Error(c, 1001, "user_id is required")
		return
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	messages, err := h.msgSvc.GetConversation(c.Request.Context(), currentUserID, targetUserID, limit, offset)
	if err != nil {
		response.Error(c, 5001, "server error")
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
		response.Error(c, 1001, "invalid request")
		return
	}

	currentUserID := c.GetString("user_id")
	if err := h.msgSvc.MarkConversationAsRead(c.Request.Context(), currentUserID, req.OtherUserID); err != nil {
		response.Error(c, 5001, "server error")
		return
	}
	response.Success(c, nil)
}
