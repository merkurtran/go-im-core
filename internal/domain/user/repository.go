package user

import "context"

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id string) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	Search(ctx context.Context, keyword string, limit, offset int) ([]*User, error)
	Update(ctx context.Context, user *User) error
	UpdatePassword(ctx context.Context, userID, hashedPassword string) error
	Delete(ctx context.Context, id string) error
}
