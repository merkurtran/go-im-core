package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/merkurtran/go-im-core/internal/domain/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
	return &mongoUserRepo{
		collection: db.Collection("users"),
	}
}

func (r *mongoUserRepo) Create(ctx context.Context, user *user.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return ErrUserAlreadyExists
		}
		return err
	}

	user.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return nil
}

func (r *mongoUserRepo) GetByID(ctx context.Context, id string) (*user.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, ErrInvalidID
	}

	var user user.User
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID, "is_deleted": false}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil

}

func (r *mongoUserRepo) GetByUsername(ctx context.Context, username string) (*user.User, error) {
	var user user.User
	err := r.collection.FindOne(ctx, bson.M{"username": username, "is_deleted": false}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *mongoUserRepo) Update(ctx context.Context, user *user.User) error {
	objectID, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		return ErrInvalidID
	}

	update := bson.M{"$set": bson.M{
		"nickname":   user.Nickname,
		"avatar":     user.Avatar,
		"status":     user.Status,
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
