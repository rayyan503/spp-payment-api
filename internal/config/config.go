package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort          string
	DBHost              string
	DBPort              string
	DBUser              string
	DBPassword          string
	DBName              string
	JWTSecretKey        string
	JWTExpirationHours  string
	MidtransServerKey   string
	MidtransClientKey   string
	MidtransEnvironment string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	return &Config{
		ServerPort:          os.Getenv("SERVER_PORT"),
		DBHost:              os.Getenv("DB_HOST"),
		DBPort:              os.Getenv("DB_PORT"),
		DBUser:              os.Getenv("DB_USER"),
		DBPassword:          os.Getenv("DB_PASSWORD"),
		DBName:              os.Getenv("DB_NAME"),
		JWTSecretKey:        os.Getenv("JWT_SECRET_KEY"),
		JWTExpirationHours:  os.Getenv("JWT_EXPIRATION_HOURS"),
		MidtransServerKey:   os.Getenv("MIDTRANS_SERVER_KEY"),
		MidtransClientKey:   os.Getenv("MIDTRANS_CLIENT_KEY"),
		MidtransEnvironment: os.Getenv("MIDTRANS_ENVIRONMENT"),
	}, nil
}
