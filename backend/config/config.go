package config

import (
	"os"
)

// Config holds all application-wide configuration settings.
type Config struct {
	PORT              string
	DB_HOST           string
	DB_USER           string
	DB_PASSWORD       string
	DB_NAME           string
	DB_PORT           string
	DB_MAX_OPEN_CONNS string
	DB_MAX_IDLE_CONNS string
	JWT_SECRET        string
	JWT_EXPIRATION    string
	CSRF_SECRET       string
	REDIS_ADDR        string
	REDIS_PASSWORD    string
	REDIS_DB          string
	ALLOWED_ORIGINS   string
	RATE_LIMIT        string
}

// NewConfig loads the environment variables into a Config struct.
func NewConfig() *Config {
	return &Config{
		PORT:              os.Getenv("PORT"),
		DB_HOST:           os.Getenv("DB_HOST"),
		DB_USER:           os.Getenv("DB_USER"),
		DB_PASSWORD:       os.Getenv("DB_PASSWORD"),
		DB_NAME:           os.Getenv("DB_NAME"),
		DB_PORT:           os.Getenv("DB_PORT"),
		DB_MAX_OPEN_CONNS: os.Getenv("DB_MAX_OPEN_CONNS"),
		DB_MAX_IDLE_CONNS: os.Getenv("DB_MAX_IDLE_CONNS"),
		JWT_SECRET:        os.Getenv("JWT_SECRET"),
		JWT_EXPIRATION:    os.Getenv("JWT_EXPIRATION"),
		CSRF_SECRET:       os.Getenv("CSRF_SECRET"),
		REDIS_ADDR:        os.Getenv("REDIS_ADDR"),
		REDIS_PASSWORD:    os.Getenv("REDIS_PASSWORD"),
		REDIS_DB:          os.Getenv("REDIS_DB"),
		ALLOWED_ORIGINS:   os.Getenv("ALLOWED_ORIGINS"),
		RATE_LIMIT:        os.Getenv("RATE_LIMIT"),
	}
}