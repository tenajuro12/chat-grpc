package config

import (
	"os"
	"strconv"
)

type Config struct {
	MongoURI           string
	MongoDB            string
	GRPCAddr           string
	JWTSecret          string
	JWTExpirationHours int
}

func LoadConfig() *Config {
	jwtExpHours, _ := strconv.Atoi(getEnv("JWT_EXPIRATION_HOURS", "24"))

	return &Config{
		MongoURI:           getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:            getEnv("MONGO_DB", "auth_service"),
		GRPCAddr:           getEnv("GRPC_ADDR", ":50051"),
		JWTSecret:          getEnv("JWT_SECRET", "your-secret-key"),
		JWTExpirationHours: jwtExpHours,
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
