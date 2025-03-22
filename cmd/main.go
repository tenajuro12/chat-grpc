package main

import (
	"context"
	"github.com/tenajuro12/chat-grpc-auth/internal/config"
	"github.com/tenajuro12/chat-grpc-auth/internal/delivery/grpc"
	"github.com/tenajuro12/chat-grpc-auth/internal/domain/repository"
	"github.com/tenajuro12/chat-grpc-auth/internal/usecase"
	"github.com/tenajuro12/chat-grpc-auth/pkg/jwt"
	pb "github.com/tenajuro12/chat-grpc-auth/pkg/pb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.LoadConfig()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoClient.Disconnect(ctx)

	userCollection := mongoClient.Database(cfg.MongoDB).Collection("users")
	userRepo := repository.NewMongoDBRepository(userCollection)

	jwtService := jwt.NewJWTService(cfg.JWTSecret, cfg.JWTExpirationHours)
	authUseCase := usecase.NewAuthUseCase(userRepo, jwtService)
	userUseCase := usecase.NewUserUseCase(userRepo)

	lis, err := net.Listen("tcp", cfg.GRPCAddr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewGRPCServer()
	authHandler := grpc.NewAuthHandler(authUseCase, userUseCase, jwtService)
	pb.RegisterAuthServiceServer(grpcServer, &authHandler)

	reflection.Register(grpcServer)

	go func() {
		log.Printf("gRPC server starting on %s", cfg.GRPCAddr)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	grpcServer.GracefulStop()
	log.Println("Server exited")

}
