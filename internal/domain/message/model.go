package message

import "time"

type Message struct {
	ID         string     `bson:"_id,omitempty" json:"message_id"`
	SenderID   string     `bson:"sender_id" json:"sender_id"`
	ReceiverID string     `bson:"receiver_id" json:"receiver_id"`
	Content    string     `bson:"content" json:"content"`
	MsgType    string     `bson:"msg_type" json:"msg_type"` // text, image, file
	Status     string     `bson:"status" json:"status"`     // sent, delivered, read
	ReadAt     *time.Time `bson:"read_at,omitempty" json:"read_at,omitempty"`
	CreatedAt  time.Time  `bson:"created_at" json:"created_at"`
}
