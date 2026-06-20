package service

import (
	"context"
	"errors"

	"github.com/merkurtran/go-im-core/internal/domain/message"
	"github.com/merkurtran/go-im-core/internal/domain/user"
	"github.com/merkurtran/go-im-core/pkg/validator"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type MessageService struct {
	repo     message.MessageRepository
	userRepo user.UserRepository
}

type SendMessageRequest struct {
	SenderID    string `json:"sender_id"`
	ReceiverID  string `json:"receiver_id"`
	Content     string `json:"content"`
	MessageType string `json:"message_type"` // text, image, file
	Status      string `json:"status"`
}
type SendMessageResponse struct {
	MessageID string `json:"message_id"`
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

	if _, err := s.userRepo.GetByID(ctx, req.SenderID); err != nil {
		return nil, ErrUserNotFound
	}
	if _, err := s.userRepo.GetByID(ctx, req.ReceiverID); err != nil {
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
