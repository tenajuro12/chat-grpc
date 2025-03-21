package repository

import (
	"context"
	"github.com/tenajuro12/chat-grpc-auth/internal/domain/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) (primitive.ObjectID, error)
	GetUserById(ctx context.Context, id primitive.ObjectID) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
}
