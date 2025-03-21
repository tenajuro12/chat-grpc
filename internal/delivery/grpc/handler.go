package grpc

import (
	"context"
	"fmt"
	"github.com/tenajuro12/chat-grpc-auth/internal/usecase"
	"github.com/tenajuro12/chat-grpc-auth/pkg/jwt"
	pb "github.com/tenajuro12/chat-grpc-auth/pkg/pb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	pb.UnimplementedAuthServiceServer
	authUseCase usecase.AuthUseCase
	userUseCase usecase.UserUseCase
	jwtService  jwt.JWTService
}

func NewAuthHandler(authUseCase usecase.AuthUseCase, userUseCase usecase.UserUseCase, jwtService jwt.JWTService) AuthHandler {
	return AuthHandler{
		authUseCase: authUseCase,
		userUseCase: userUseCase,
		jwtService:  jwtService,
	}
}

func (h *AuthHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.AuthResponse, error) {
	user, err := h.userUseCase.CreateUser(ctx, req.Username, req.Email, req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	token, err := h.jwtService.GenerateToken(user.ID.Hex())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate token: %v", err)
	}

	return &pb.AuthResponse{
		Token:    token,
		UserId:   user.ID.Hex(),
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

func (h *AuthHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.AuthResponse, error) {
	token, err := h.authUseCase.Login(ctx, req.UsernameOrEmail, req.Password)
	if err != nil {
		return nil, fmt.Errorf("error with loggin in", err)
	}
	userID, err := h.jwtService.ValidateToken(token)
	if err != nil {
		return nil, fmt.Errorf("error with validating token", err)
	}
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "invalid user ID format: %v", err)
	}

	user, err := h.userUseCase.GetUserByID(ctx, objectID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}
	return &pb.AuthResponse{
		Token:    token,
		UserId:   userID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

func (h *AuthHandler) Validate(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	userID, err := h.authUseCase.ValidateToken(req.Token)
	if err != nil {
		return &pb.ValidateResponse{
			Valid:  false,
			UserId: "",
		}, nil
	}

	return &pb.ValidateResponse{
		Valid:  true,
		UserId: userID,
	}, nil
}
