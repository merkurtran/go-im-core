package message

import "context"

type MessageRepository interface {
	Create(ctx context.Context, message *Message) error
	GetByID(ctx context.Context, id string) (*Message, error)

	// Get between two users conversation with limit and offset
	GetByConversation(ctx context.Context, userID1, userID2 string, limit, offset int) ([]*Message, error)

	// Get unread count of a user
	GetUnreadCount(ctx context.Context, userID string) (int, error)

	// Mark a message as read
	MarkAsRead(ctx context.Context, messageID string) error

	MarkConversationAsRead(ctx context.Context, currentUserID, otherUserID string) error

	Update(ctx context.Context, message *Message) error
	Delete(ctx context.Context, id string) error
}
