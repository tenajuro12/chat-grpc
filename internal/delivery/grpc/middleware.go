package grpc

import (
	"context"
	"github.com/tenajuro12/chat-grpc-auth/pkg/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthInterceptor struct {
	jwtService jwt.JWTService
}

func NewAuthInterceptor(jwtService jwt.JWTService) *AuthInterceptor {
	return &AuthInterceptor{
		jwtService: jwtService,
	}
}

func (interceptor *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if info.FullMethod == "/auth.AuthService/Login" || info.FullMethod == "/auth.AuthService/Register" || info.FullMethod == "/auth.AuthService/Validate" {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
		}

		values := md["authorization"]
		if len(values) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
		}

		userID, err := interceptor.jwtService.ValidateToken(values[0])
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
		}

		newCtx := context.WithValue(ctx, "user_id", userID)
		return handler(newCtx, req)
	}
}

func NewGRPCServer() *grpc.Server {
	return grpc.NewServer()
}
