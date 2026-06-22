package mongo

import (
	"time"

	"github.com/merkurtran/go-im-core/internal/domain/message"
	"github.com/merkurtran/go-im-core/internal/domain/user"
)

// ===== User DTO =====

type userDTO struct {
	ID        string    `bson:"_id,omitempty"`
	Username  string    `bson:"username"`
	Password  string    `bson:"password"`
	Nickname  string    `bson:"nickname"`
	Avatar    string    `bson:"avatar"`
	Status    string    `bson:"status"`
	IsDeleted bool      `bson:"is_deleted"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

func toUserDTO(u *user.User) *userDTO {
	if u == nil {
		return nil
	}
	return &userDTO{
		ID:        u.ID,
		Username:  u.Username,
		Password:  u.Password,
		Nickname:  u.Nickname,
		Avatar:    u.Avatar,
		Status:    u.Status,
		IsDeleted: u.IsDeleted,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func toUserModel(d *userDTO) *user.User {
	if d == nil {
		return nil
	}
	return &user.User{
		ID:        d.ID,
		Username:  d.Username,
		Password:  d.Password,
		Nickname:  d.Nickname,
		Avatar:    d.Avatar,
		Status:    d.Status,
		IsDeleted: d.IsDeleted,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}

// ===== Message DTO =====

type messageDTO struct {
	ID         string     `bson:"_id,omitempty"`
	SenderID   string     `bson:"sender_id"`
	ReceiverID string     `bson:"receiver_id"`
	Content    string     `bson:"content"`
	MsgType    string     `bson:"msg_type"`
	Status     string     `bson:"status"`
	ReadAt     *time.Time `bson:"read_at,omitempty"`
	IsDeleted  bool       `bson:"is_deleted"`
	CreatedAt  time.Time  `bson:"created_at"`
}

func toMessageDTO(m *message.Message) *messageDTO {
	if m == nil {
		return nil
	}
	return &messageDTO{
		ID:         m.ID,
		SenderID:   m.SenderID,
		ReceiverID: m.ReceiverID,
		Content:    m.Content,
		MsgType:    m.MsgType,
		Status:     m.Status,
		ReadAt:     m.ReadAt,
		IsDeleted:  m.IsDeleted,
		CreatedAt:  m.CreatedAt,
	}
}

func toMessageModel(d *messageDTO) *message.Message {
	if d == nil {
		return nil
	}
	return &message.Message{
		ID:         d.ID,
		SenderID:   d.SenderID,
		ReceiverID: d.ReceiverID,
		Content:    d.Content,
		MsgType:    d.MsgType,
		Status:     d.Status,
		ReadAt:     d.ReadAt,
		IsDeleted:  d.IsDeleted,
		CreatedAt:  d.CreatedAt,
	}
}

func toMessageModels(ds []*messageDTO) []*message.Message {
	if len(ds) == 0 {
		return nil
	}
	out := make([]*message.Message, 0, len(ds))
	for _, d := range ds {
		if m := toMessageModel(d); m != nil {
			out = append(out, m)
		}
	}
	return out
}
