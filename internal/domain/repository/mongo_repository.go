package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/tenajuro12/chat-grpc-auth/internal/domain/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoDBRepository struct {
	collection *mongo.Collection
}

func NewMongoDBRepository(collection *mongo.Collection) UserRepository {
	return mongoDBRepository{collection: collection}
}

func (m mongoDBRepository) CreateUser(ctx context.Context, user *models.User) (primitive.ObjectID, error) {
	result, err := m.collection.InsertOne(ctx, &user)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("failed to insert user: %w", err)
	}
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, fmt.Errorf("failed to cast InsertedID to ObjectID: %v", result.InsertedID)
	}
	return oid, err
}

func (m mongoDBRepository) GetUserById(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	var user models.User
	err := m.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &user, err
}

func (m mongoDBRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := m.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &user, err
}

func (m mongoDBRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := m.collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &user, err
}

func (m mongoDBRepository) UpdateUser(ctx context.Context, user *models.User) error {
	_, err := m.collection.ReplaceOne(ctx, bson.M{"_id": user.ID}, user)
	if err != nil {
		return fmt.Errorf("Failed to update user", err)
	}
	return err
}
