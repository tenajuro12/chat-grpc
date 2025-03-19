package main

import (
	"github.com/tenajuro12/chat-grpc/internal/auth"
	"github.com/tenajuro12/chat-grpc/internal/storage"
	pb "github.com/tenajuro12/chat-grpc/pkg/auth_v1"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	db := storage.InitMongoDB()
	authService := auth.NewAuthService(db)
	listener, err := net.Listen("tcp", ":50501")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, authService)
	log.Println("✅ gRPC Server is running on port 50051...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
