package message

import (
	"errors"
	"time"
)

var (
	ErrMessageNotFound  = errors.New("message not found")
	ErrNotMessageSender = errors.New("you are not the sender of this message")
	ErrRecallTimeout    = errors.New("recall time has expired")
)

type Message struct {
	ID         string     `json:"message_id"`
	SenderID   string     `json:"sender_id"`
	ReceiverID string     `json:"receiver_id"`
	Content    string     `json:"content"`
	MsgType    string     `json:"msg_type"` // text, image, file
	Status     string     `json:"status"`   // sent, delivered, read
	ReadAt     *time.Time `json:"read_at,omitempty"`
	IsDeleted  bool       `json:"is_deleted"`
	IsRecalled bool       `json:"is_recalled"`
	CreatedAt  time.Time  `json:"created_at"`
}
