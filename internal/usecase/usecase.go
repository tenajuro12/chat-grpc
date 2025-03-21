package usecase

import (
	"context"
	"errors"
	"github.com/tenajuro12/chat-grpc-auth/internal/domain/models"
	"github.com/tenajuro12/chat-grpc-auth/internal/domain/repository"
	"github.com/tenajuro12/chat-grpc-auth/pkg/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type UserUseCase interface {
	CreateUser(ctx context.Context, username, email, password string) (*models.User, error)
	GetUserByID(ctx context.Context, id primitive.ObjectID) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
}

type userUseCase struct {
	userRepo repository.UserRepository
}

func NewUserUseCase(userRepo repository.UserRepository) UserUseCase {
	return &userUseCase{userRepo: userRepo}
}

func (uc *userUseCase) CreateUser(ctx context.Context, username, email, password string) (*models.User, error) {
	existingUser, err := uc.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	existingUser, err = uc.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("user with this username already exists")
	}

	passwordHash, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	user := &models.User{
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	id, err := uc.userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	user.ID = id
	return user, nil
}

func (uc *userUseCase) GetUserByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	return uc.userRepo.GetUserById(ctx, id)
}

func (uc *userUseCase) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return uc.userRepo.GetUserByEmail(ctx, email)
}

func (uc *userUseCase) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	return uc.userRepo.GetUserByUsername(ctx, username)
}
