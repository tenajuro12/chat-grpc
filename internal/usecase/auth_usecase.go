package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/tenajuro12/chat-grpc-auth/internal/domain/repository"
	"github.com/tenajuro12/chat-grpc-auth/pkg/jwt"
	"github.com/tenajuro12/chat-grpc-auth/pkg/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthUseCase interface {
	Login(ctx context.Context, usernameOrEmail, password string) (string, error)
	ValidateToken(token string) (string, error)
}

type authUseCase struct {
	userRepo   repository.UserRepository
	jwtService jwt.JWTService
}

func NewAuthUseCase(userRepo repository.UserRepository) AuthUseCase {
	return authUseCase{userRepo: userRepo}
}

func (a authUseCase) Login(ctx context.Context, usernameOrEmail, password string) (string, error) {
	user, err := a.userRepo.GetUserByEmail(ctx, usernameOrEmail)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return "", fmt.Errorf("failed to get user by email: %w", err)
	}

	if user == nil {
		user, err = a.userRepo.GetUserByUsername(ctx, usernameOrEmail)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return "", errors.New("user not found") // Обычная ошибка, без %w
			}
			return "", fmt.Errorf("failed to get user by username: %w", err)
		}
	}

	if !utils.CheckPassword(password, user.PasswordHash) {
		return "", errors.New("invalid password")
	}

	token, err := a.jwtService.GenerateToken(user.ID.Hex())
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}

func (a authUseCase) ValidateToken(token string) (string, error) {
	return a.jwtService.ValidateToken(token)
}
