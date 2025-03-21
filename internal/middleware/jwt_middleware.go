package middleware

import (
	"context"
	"github.com/tenajuro12/chat-grpc-auth/pkg/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
)

func JWTAuthMiddleware(jwtService jwt.JWTService) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if info.FullMethod == "/auth.AuthService/Login" || info.FullMethod == "/auth.AuthService/Register" || info.FullMethod == "/auth.AuthService/Validate" {
			return handler(ctx, req)
		}

		userID, err := extractUserID(ctx, jwtService)
		if err != nil {
			return nil, err
		}

		newCtx := context.WithValue(ctx, "user_id", userID)
		return handler(newCtx, req)
	}
}

func extractUserID(ctx context.Context, jwtService jwt.JWTService) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return "", status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	authHeader := values[0]
	fields := strings.Fields(authHeader)
	if len(fields) < 2 {
		return "", status.Errorf(codes.Unauthenticated, "invalid authorization format")
	}

	authType := strings.ToLower(fields[0])
	if authType != "bearer" {
		return "", status.Errorf(codes.Unauthenticated, "unsupported authorization type")
	}

	accessToken := fields[1]
	userID, err := jwtService.ValidateToken(accessToken)
	if err != nil {
		return "", status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	return userID, nil
}
