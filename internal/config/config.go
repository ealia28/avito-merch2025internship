// Package config загружает конфигурацию сервиса из переменных окружения.
package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

const minJWTSecretLen = 32

// Server — настройки HTTP-сервера: порт и таймауты.
type Server struct {
	Port              string        `env:"PORT" envDefault:"8080"`
	ReadHeaderTimeout time.Duration `env:"READ_HEADER_TIMEOUT" envDefault:"5s"`
	ReadTimeout       time.Duration `env:"READ_TIMEOUT" envDefault:"10s"`
	WriteTimeout      time.Duration `env:"WRITE_TIMEOUT" envDefault:"10s"`
	IdleTimeout       time.Duration `env:"IDLE_TIMEOUT" envDefault:"60s"`
	ShutdownTimeout   time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"10s"`
}

// Database — параметры подключения к PostgreSQL.
type Database struct {
	Host     string `env:"HOST" envDefault:"localhost"`
	Port     string `env:"PORT" envDefault:"5432"`
	User     string `env:"USER,required"`
	Password string `env:"PASSWORD,required"`
	Name     string `env:"NAME,required"`
}

// JWT — настройки выпуска и проверки токенов.
type JWT struct {
	Secret string        `env:"SECRET,required"`
	TTL    time.Duration `env:"TTL" envDefault:"24h"`
}

// Config — полная конфигурация сервиса.
type Config struct {
	Server   Server   `envPrefix:"SERVER_"`
	Database Database `envPrefix:"DATABASE_"`
	JWT      JWT      `envPrefix:"JWT_"`
}

// Load читает конфигурацию из переменных окружения и проверяет её.
// Возвращает ошибку, если не заданы обязательные переменные
// или JWT-секрет короче minJWTSecretLen байт.
func Load() (*Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("parse env: %w", err)
	}
	if len(cfg.JWT.Secret) < minJWTSecretLen {
		return nil, fmt.Errorf("jwt secret must be at least %d bytes long", minJWTSecretLen)
	}
	return &cfg, nil
}
