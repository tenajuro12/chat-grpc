package auth

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/tenajuro12/chat-grpc/pkg/auth_v1"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

const secret_key = "Supersecretkey"

type AuthService struct {
	auth_v1.UnimplementedAuthServiceServer
	userCollection *mongo.Collection
}

func NewAuthService(db *mongo.Database) *AuthService {
	return &AuthService{userCollection: db.Collection("users")}
}

func (s *AuthService) Register(ctx context.Context, req *auth_v1.RegisterRequest) (*auth_v1.AuthResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Error with hashins")
	}
	existingUser := s.userCollection.FindOne(ctx, bson.M{"username": req.Username})
	if existingUser.Err() == nil {
		return nil, errors.New("user already exists")
	}

	_, err = s.userCollection.InsertOne(ctx, bson.M{
		"username": req.Username,
		"password": string(hashedPassword),
	})
	if err != nil {
		return nil, err
	}
	token, err := generateJWT(req.Username)
	if err != nil {
		return nil, err
	}
	log.Println("✅ User logged in:", req.Username)
	return &auth_v1.AuthResponse{Token: token}, nil

}

func (s *AuthService) Login(ctx context.Context, req *auth_v1.LoginRequest) (*auth_v1.AuthResponse, error) {
	var user struct {
		Username string `bson:"username"`
		Password string `bson:"password"`
	}

	err := s.userCollection.FindOne(ctx, bson.M{
		"username": req.Username,
	}).Decode(&user)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid username or password")
	}
	token, err := generateJWT(user.Username)
	if err != nil {
		return nil, err
	}

	log.Println("✅ User logged in:", req.Username)
	return &auth_v1.AuthResponse{Token: token}, nil
}

func generateJWT(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})
	return token.SignedString([]byte(secret_key))
}
