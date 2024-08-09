package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerHost             string
	ServerPort             string
	DBName                 string
	DBConnectionString     string
	JWTSecret              string
	JWTExpirationInSeconds int64
}

var Envs = initConfig()

func initConfig() Config {
	if err := godotenv.Load("example.env"); err != nil {
		log.Panicln(err)
	}

	return Config{
		ServerHost:             getEnv("PUBLIC_HOST", "http://localhost"),
		ServerPort:             getEnv("PORT", "8080"),
		DBConnectionString:     getEnv("DB_CONNECTION_STRING", ""),
		DBName:                 getEnv("DB_NAME", "taskdb"),
		JWTSecret:              getEnv("JWT_SECRET", "not-so-secret-now-is-it?"),
		JWTExpirationInSeconds: getEnvAsInt("JWT_EXPIRATION_IN_SECONDS", 60*24),
	}
}

// Gets the env by key or fallbacks
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func getEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}

		return i
	}

	return fallback
}
