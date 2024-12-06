package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBDriver      string
	DBSource      string
	ServerAddress string
	JWTSecret     string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	config := &Config{
		DBDriver:      getEnv("DB_DRIVER", "mysql"),
		DBSource:      getEnv("DB_SOURCE", "root:password@tcp(localhost:3306)/attendance_db?parseTime=true"),
		ServerAddress: getEnv("SERVER_ADDRESS", ":8080"),
		JWTSecret:     getEnv("JWT_SECRET", "your-secret-key"),
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
