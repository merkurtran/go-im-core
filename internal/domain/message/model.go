package message

import "time"

type Message struct {
	ID         string     `json:"message_id"`
	SenderID   string     `json:"sender_id"`
	ReceiverID string     `json:"receiver_id"`
	Content    string     `json:"content"`
	MsgType    string     `json:"msg_type"` // text, image, file
	Status     string     `json:"status"`     // sent, delivered, read
	ReadAt     *time.Time `json:"read_at,omitempty"`
	IsDeleted  bool       `json:"is_deleted"`
	CreatedAt  time.Time  `json:"created_at"`
}
