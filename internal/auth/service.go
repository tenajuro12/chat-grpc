package auth

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
	pb "github.com/tenajuro12/chat-grpc/pkg/auth_v1"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

const secret_key = "Supersecretkey"

type AuthService struct {
	pb.Unim
}
