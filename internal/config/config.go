package config

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Cookie   CookieConfig
	Modules  CustomModule
}

type CustomModule struct {
	Notify bool
}

type ServerConfig struct {
	Host string
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type JWTConfig struct {
	Secret        string
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration
}

type CookieConfig struct {
	Domain   string
	Secure   string
	SameSite http.SameSite
}

func Load() (*Config, error) {
	// TODO: Load from .env file
	accessExpiry, err := time.ParseDuration(getEnv("JWT_ACCESS_EXPIRY", "15m"))
	if err != nil {
		return nil, fmt.Errorf("JWT_ACCESS_EXPIRY invalid: %w", err)
	}

	refreshExpiry, err := time.ParseDuration(getEnv("JWT_REFRESH_EXPIRY", "7d"))
	if err != nil {
		return nil, fmt.Errorf("JWT_REFRESH_EXPIRY invalid: %w", err)
	}

	sameSite, err := strconv.ParseInt(getEnv("COOKIE_SAMESITE", "1"), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("COOKIE_SAME_SITE invalid: %w", err)
	}

	notifyModule, err := strconv.ParseBool(getEnv("NOTIFY_MODULE", "false"))
	if err != nil {
		notifyModule = false
		log.Println("NOTIFY_MODULE is invalid, continue with false")
	}

	return &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "127.0.0.1"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "nexa_tracker"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:        getEnv("JWT_SECRET", ""),
			AccessExpiry:  accessExpiry,
			RefreshExpiry: refreshExpiry,
		},
		Cookie: CookieConfig{
			Domain:   getEnv("COOKIE_DOMAIN", ""),
			Secure:   getEnv("COOKIE_SECURE", "false"),
			SameSite: http.SameSite(sameSite),
		},
		Modules: CustomModule{
			Notify: notifyModule,
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
