package user

import "time"

type User struct {
	ID        string    `bson:"_id,omitempty" json:"user_id"`
	Username  string    `bson:"username" json:"username"`
	Password  string    `bson:"password" json:"-"`
	Nickname  string    `bson:"nickname" json:"nickname"`
	Avatar    string    `bson:"avatar" json:"avatar"`
	Status    string    `bson:"status" json:"status"`
	IsDeleted bool      `bson:"is_deleted" json:"is_deleted"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}
