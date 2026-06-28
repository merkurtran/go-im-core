package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/merkurtran/go-im-core/internal/domain/message"
	"github.com/merkurtran/go-im-core/internal/domain/user"
	"github.com/merkurtran/go-im-core/pkg/validator"
)

// var (
// 	ErrUserNotFound = errors.New("user not found")
// )

type MessageService struct {
	repo     message.MessageRepository
	userRepo user.UserRepository
}

type SendMessageRequest struct {
	SenderID    string `json:"-"` // 服务端从 auth context 获取，客户端无需传入
	ReceiverID  string `json:"receiver_id"`
	Content     string `json:"content"`
	MessageType string `json:"message_type"` // text, image, file
}

type SendMessageResponse struct {
	MessageID string `json:"message_id"`
}

type MarkReadRequest struct {
	OtherUserID string `json:"other_user_id"`
}

func NewMessageService(repo message.MessageRepository, userRepo user.UserRepository) *MessageService {
	return &MessageService{
		repo:     repo,
		userRepo: userRepo,
	}
}

func (s *MessageService) SendMessage(ctx context.Context, req *SendMessageRequest) (*SendMessageResponse, error) {
	if err := validator.ValidateMessage(req.SenderID, req.ReceiverID, req.Content, req.MessageType); err != nil {
		return nil, err
	}

	sender, err := s.userRepo.GetByID(ctx, req.SenderID)
	if err != nil {
		return nil, err
	}
	if sender == nil {
		return nil, ErrUserNotFound
	}

	receiver, err := s.userRepo.GetByID(ctx, req.ReceiverID)
	if err != nil {
		return nil, err
	}
	if receiver == nil {
		return nil, ErrUserNotFound
	}

	msg := &message.Message{
		SenderID:   req.SenderID,
		ReceiverID: req.ReceiverID,
		Content:    req.Content,
		MsgType:    req.MessageType,
		Status:     "sent",
	}
	if err := s.repo.Create(ctx, msg); err != nil {
		return nil, err
	}

	return &SendMessageResponse{
		MessageID: msg.ID,
	}, nil
}

func (s *MessageService) GetConversation(ctx context.Context, currentUserID, targetUserID string, limit, offset int) ([]*message.Message, error) {
	return s.repo.GetByConversation(ctx, currentUserID, targetUserID, limit, offset)
}

func (s *MessageService) GetUnreadCount(ctx context.Context, userID string) (int, error) {
	return s.repo.GetUnreadCount(ctx, userID)
}

func (s *MessageService) MarkConversationAsRead(ctx context.Context, currentUserID, otherUserID string) error {
	return s.repo.MarkConversationAsRead(ctx, currentUserID, otherUserID)
}

func (s *MessageService) DeleteMessage(ctx context.Context, messageID string) error {
	return s.repo.Delete(ctx, messageID)
}

func (s *MessageService) RecallMessage(ctx context.Context, userID, messageID string) error {
	msg, err := s.repo.GetByID(ctx, messageID)
	if err != nil {
		return err
	}
	if msg == nil {
		return message.ErrMessageNotFound
	}
	if msg.SenderID != userID {
		return message.ErrNotMessageSender
	}
	if time.Since(msg.CreatedAt) > 2*time.Minute {
		return message.ErrRecallTimeout
	}
	msg.IsRecalled = true
	return s.repo.Update(ctx, msg)
}

// GetMessageByID 用于 WebSocket 业务层获取消息详情
func (s *MessageService) GetMessageByID(ctx context.Context, id string) (*message.Message, error) {
	return s.repo.GetByID(ctx, id)
}

// OnUserConnected 用户上线时的回调（更新状态）
func (s *MessageService) OnUserConnected(ctx context.Context, userID string) {
	// 业务层可以在这里扩展：推送离线消息、更新已读状态等
	slog.Info("user connected", "user_id", userID)
}

// OnUserDisconnected 用户离线时的回调
func (s *MessageService) OnUserDisconnected(ctx context.Context, userID string) {
	slog.Info("user disconnected", "user_id", userID)
}
