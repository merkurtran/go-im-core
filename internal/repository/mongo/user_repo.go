package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/merkurtran/go-im-core/internal/domain/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidID         = errors.New("invalid user id format")
)

type mongoUserRepo struct {
	collection *mongo.Collection
}

func NewUserRepo(db *mongo.Database) user.UserRepository {
	repo := &mongoUserRepo{
		collection: db.Collection("users"),
	}
	// 初始化索引：username 唯一
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	repo.collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "username", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	return repo
}

func (r *mongoUserRepo) Create(ctx context.Context, u *user.User) error {
	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now

	dto := toUserDTO(u)
	result, err := r.collection.InsertOne(ctx, dto)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return ErrUserAlreadyExists
		}
		return err
	}

	u.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return nil
}

func (r *mongoUserRepo) GetByID(ctx context.Context, id string) (*user.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, ErrInvalidID
	}

	var dto userDTO
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID, "is_deleted": false}).Decode(&dto)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return toUserModel(&dto), nil
}

func (r *mongoUserRepo) GetByUsername(ctx context.Context, username string) (*user.User, error) {
	var dto userDTO
	err := r.collection.FindOne(ctx, bson.M{"username": username, "is_deleted": false}).Decode(&dto)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return toUserModel(&dto), nil
}

func (r *mongoUserRepo) Update(ctx context.Context, u *user.User) error {
	objectID, err := primitive.ObjectIDFromHex(u.ID)
	if err != nil {
		return ErrInvalidID
	}

	update := bson.M{"$set": bson.M{
		"nickname":   u.Nickname,
		"avatar":     u.Avatar,
		"status":     u.Status,
		"updated_at": time.Now(),
	}}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (r *mongoUserRepo) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrInvalidID
	}
	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": bson.M{"is_deleted": true, "updated_at": time.Now()}})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return ErrUserNotFound
	}
	return nil
}
