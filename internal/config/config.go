package config

import (
	"github.com/joho/godotenv"
	"os"
	"rest-api-tutorial/pkg/logging"
	"strconv"
)

type Config struct {
	IsDebug    bool
	Listen     Listen
	PostgreSQL PostgreSQL
}

type Listen struct {
	Type   string
	BindIP string
	Port   string
}

type PostgreSQL struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

type User struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

func LoadConfigEnv() *Config {
	logger := logging.GetLogger()
	if err := godotenv.Load(); err != nil {
		logger.Warn("Warning: .env file not found, using environment variables")
	}

	return &Config{
		PostgreSQL: PostgreSQL{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			Username: getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			Database: getEnv("DB_NAME", "db"),
		},
		Listen: Listen{
			Type:   getEnv("LISTEN_TYPE", "http"),
			BindIP: getEnv("BIND_IP", "0.0.0.0"),
			Port:   getEnv("APP_PORT", "8080"),
		},
		IsDebug: getEnvAsBool("DEBUG", false),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	value := getEnv(key, "")
	if v, err := strconv.ParseBool(value); err == nil {
		return v
	}
	return defaultValue
}
