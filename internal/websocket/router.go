package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/merkurtran/go-im-core/internal/service"
)

// DefaultMessageRouter 实现 MessageRouter 接口，将 WebSocket 消息路由到业务层。
type DefaultMessageRouter struct {
	msgService  *service.MessageService
	userService *service.UserService
	hub         *Hub
}

func NewDefaultMessageRouter(
	msgSvc *service.MessageService,
	userSvc *service.UserService,
	hub *Hub,
) *DefaultMessageRouter {
	return &DefaultMessageRouter{
		msgService:  msgSvc,
		userService: userSvc,
		hub:         hub,
	}
}

func (r *DefaultMessageRouter) OnMessage(ctx context.Context, senderID string, msg *WSMessage) (*WSMessage, error) {
	switch msg.Event {
	case "message":
		return r.handleMessage(ctx, senderID, msg.Payload)
	case "ack":
		return r.handleAck(ctx, senderID, msg.Payload)
	case "ping":
		return &WSMessage{Event: "pong"}, nil
	default:
		return nil, fmt.Errorf("unknown event: %s", msg.Event)
	}
}

func (r *DefaultMessageRouter) handleMessage(ctx context.Context, senderID string, payload json.RawMessage) (*WSMessage, error) {
	var p MessagePayload
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil, err
	}

	req := &service.SendMessageRequest{
		SenderID:    senderID,
		ReceiverID:  p.ReceiverID,
		Content:     p.Content,
		MessageType: p.MessageType,
	}
	resp, err := r.msgService.SendMessage(ctx, req)
	if err != nil {
		return nil, err
	}

	// 推送消息给接收方
	msgBody := &WSMessage{
		Event: "message",
		Payload: mustJSON(map[string]any{
			"message_id":  resp.MessageID,
			"sender_id":   senderID,
			"receiver_id": p.ReceiverID,
			"content":     p.Content,
			"msg_type":    p.MessageType,
			"status":      "sent",
		}),
	}
	msgBytes, _ := json.Marshal(msgBody)
	r.hub.SendToUserAsync(p.ReceiverID, msgBytes)

	// 给发送方回执
	return &WSMessage{
		Event: "ack",
		Payload: mustJSON(map[string]any{
			"message_id": resp.MessageID,
			"status":     "sent",
		}),
	}, nil
}

func (r *DefaultMessageRouter) handleAck(ctx context.Context, senderID string, payload json.RawMessage) (*WSMessage, error) {
	var p AckPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil, err
	}

	if p.Status == "read" && p.OtherUserID != "" {
		if err := r.msgService.MarkConversationAsRead(ctx, senderID, p.OtherUserID); err != nil {
			slog.Warn("mark conversation as read failed", "error", err)
		}
	}
	return nil, nil
}

func mustJSON(v any) json.RawMessage {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}
