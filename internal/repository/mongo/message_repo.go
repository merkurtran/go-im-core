package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/merkurtran/go-im-core/internal/domain/message"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrMessageNotFound = errors.New("message not found")
)

type mongoMessageRepo struct {
	collection *mongo.Collection
}

func NewMessageRepo(db *mongo.Database) message.MessageRepository {
	return &mongoMessageRepo{
		collection: db.Collection("messages"),
	}
}

func (r *mongoMessageRepo) Create(ctx context.Context, message *message.Message) error {
	message.CreatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, message)
	if err != nil {
		return err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		message.ID = oid.Hex()
	}
	return nil
}

func (r *mongoMessageRepo) GetByID(ctx context.Context, id string) (*message.Message, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, ErrInvalidID
	}

	var msg message.Message
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID, "is_deleted": false}).Decode(&msg)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrMessageNotFound
		}
		return nil, err
	}

	return &msg, nil
}

func (r *mongoMessageRepo) GetByConversation(ctx context.Context, userID1, userID2 string, limit, offset int) ([]*message.Message, error) {
	var messages []*message.Message

	skip := int64(offset)
	limited := int64(limit)
	opts := &options.FindOptions{
		Sort:  bson.D{{Key: "created_at", Value: -1}},
		Skip:  &skip,
		Limit: &limited,
	}

	cursor, err := r.collection.Find(ctx, bson.M{"$or": []bson.M{
		{"sender_id": userID1, "receiver_id": userID2},
		{"sender_id": userID2, "receiver_id": userID1},
	}, "is_deleted": false}, opts)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &messages)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (r *mongoMessageRepo) GetUnreadCount(ctx context.Context, userID string) (int, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{
		"receiver_id": userID,
		"status":      bson.M{"$ne": "read"},
		"is_deleted":  false,
	})
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *mongoMessageRepo) MarkAsRead(ctx context.Context, messageID string) error {
	objectID, err := primitive.ObjectIDFromHex(messageID)
	if err != nil {
		return ErrInvalidID
	}

	now := time.Now()
	result, err := r.collection.UpdateOne(ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": bson.M{"status": "read", "read_at": now}},
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return ErrMessageNotFound
	}
	return nil
}

func (r *mongoMessageRepo) MarkConversationAsRead(ctx context.Context, currentUserID, otherUserID string) error {
	_, err := r.collection.UpdateMany(ctx,
		bson.M{"sender_id": otherUserID, "receiver_id": currentUserID, "status": bson.M{"$ne": "read"}},
		bson.M{"$set": bson.M{"status": "read", "read_at": time.Now()}},
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *mongoMessageRepo) Update(ctx context.Context, message *message.Message) error {
	objectID, err := primitive.ObjectIDFromHex(message.ID)
	if err != nil {
		return ErrInvalidID
	}

	result, err := r.collection.UpdateOne(ctx,
		bson.M{"_id": objectID, "is_deleted": false},
		bson.M{"$set": bson.M{
			"content":  message.Content,
			"msg_type": message.MsgType,
			"status":   message.Status,
		}},
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return ErrMessageNotFound
	}
	return nil
}

func (r *mongoMessageRepo) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrInvalidID
	}

	result, err := r.collection.UpdateOne(ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": bson.M{"is_deleted": true}},
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return ErrMessageNotFound
	}
	return nil
}
