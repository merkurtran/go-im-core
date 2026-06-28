package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/merkurtran/go-im-core/internal/domain/message"
	"github.com/merkurtran/go-im-core/internal/response"
	"github.com/merkurtran/go-im-core/internal/service"
)

// swaggerDocTypes 用于 swaggo 类型解析，无实际运行效果
var _ = (*message.Message)(nil)

type MessageHandler struct {
	msgSvc *service.MessageService
}

func NewMessageHandler(msgSvc *service.MessageService) *MessageHandler {
	return &MessageHandler{msgSvc: msgSvc}
}

// SendMessage  POST /messages
//
//	@Summary 发送信息
//	@Description 发送信息
//	@Tags 消息
//	@Accept json
//	@Produce json
//	@Param request body service.SendMessageRequest true "发送的信息请求"
//	@Success 200 {object} response.Response{data=service.SendMessageResponse} "发送成功"
//	@Failure 400 {object} response.Response "请求参数错误"
//	@Failure 404 {object} response.Response "用户不存在"
//	@Failure 500 {object} response.Response "服务器错误"
//	@Security Bearer
//	@Router /messages [post]
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

// GetConversation 获取聊天记录
//
//	@Summary 获取聊天记录
//	@Description 获取与指定用户的聊天消息列表
//	@Tags 消息
//	@Produce json
//	@Param user_id query string true "对方的用户ID"
//	@Param limit query int false "每页数量，默认20，最大100"
//	@Param offset query int false "偏移量"
//	@Success 200 {object} response.Response{data=[]message.Message} "成功"
//	@Failure 400 {object} response.Response "请求参数错误"
//	@Failure 500 {object} response.Response "服务器错误"
//	@Security Bearer
//	@Router /messages [get]
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

// MarkRead 标记已读
//
//	@Summary 标记会话已读
//	@Description 标记与指定用户的对话为已读
//	@Tags 消息
//	@Accept json
//	@Produce json
//	@Param request body service.MarkReadRequest true "标记已读请求"
//	@Success 200 {object} response.Response "标记成功"
//	@Failure 400 {object} response.Response "请求参数错误"
//	@Failure 500 {object} response.Response "服务器错误"
//	@Security Bearer
//	@Router /messages/read [patch]
func (h *MessageHandler) MarkRead(c *gin.Context) {
	var req service.MarkReadRequest
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

// RecallMessage 撤回消息
//
//	@Summary 撤回消息
//	@Description 撤回自己发送的消息（发送后2分钟内可撤回）
//	@Tags 消息
//	@Produce json
//	@Param message_id path string true "消息ID"
//	@Success 200 {object} response.Response "撤回成功"
//	@Failure 400 {object} response.Response "请求参数错误"
//	@Failure 403 {object} response.Response "无权撤回/超过撤回时间"
//	@Failure 404 {object} response.Response "消息不存在"
//	@Failure 500 {object} response.Response "服务器错误"
//	@Security Bearer
//	@Router /messages/{message_id}/recall [patch]
func (h *MessageHandler) RecallMessage(c *gin.Context) {
	messageID := c.Param("message_id")
	if messageID == "" {
		response.Error(c, http.StatusBadRequest, 1001, "message_id is required")
		return
	}
	userID := c.GetString("user_id")
	if userID == "" {
		response.Error(c, http.StatusUnauthorized, 1002, "unauthorized")
		return
	}
	if err := h.msgSvc.RecallMessage(c.Request.Context(), userID, messageID); err != nil {
		switch {
		case errors.Is(err, message.ErrMessageNotFound):
			response.Error(c, http.StatusNotFound, 2001, "message not found")
		case errors.Is(err, message.ErrNotMessageSender):
			response.Error(c, http.StatusForbidden, 2002, "you are not the sender")
		case errors.Is(err, message.ErrRecallTimeout):
			response.Error(c, http.StatusForbidden, 2003, "recall time has expired")
		default:
			response.Error(c, http.StatusInternalServerError, 5001, "server error")
		}
		return
	}
	response.Success(c, nil)
}
