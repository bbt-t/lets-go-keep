package config

import (
	"time"

	"github.com/caarlos0/env/v8"
	log "github.com/sirupsen/logrus"
)

// ServerConfig struct for server config.
type ServerConfig struct {
	RunAddress      string `env:"SERVER_PORT" envDefault:":3200"`
	DBConnectionURL string `env:"DATABASE_DSN"`
	FilesDirectory  string `env:"FILE_STORAGE_PATH" envDefault:"files"`
	Auth            AuthConfig
}

// AuthConfig auth settings.
type AuthConfig struct {
	SecretJWT      string `env:"SECRET_JWT" envDefault:"HERE_MUST_BE_SECRET_KEY"`
	ExpirationTime int64
}

// NewServerConfig gets server config.
func NewServerConfig() ServerConfig {
	var (
		cfg  ServerConfig
		auth AuthConfig
	)

	if err := env.Parse(&cfg); err != nil {
		log.Printf("%+v\n", err)
	}
	if err := env.Parse(&auth); err != nil {
		log.Printf("%+v\n", err)
	}

	cfg.Auth = AuthConfig{
		SecretJWT:      auth.SecretJWT,
		ExpirationTime: time.Now().Add(1 * time.Hour).Unix(),
	}
	log.Infoln("Config loaded")

	return cfg
}

// NewConfigForTests config fot tests.
func NewConfigForTests() ServerConfig {
	return ServerConfig{
		RunAddress:      ":3200",
		DBConnectionURL: "",
		FilesDirectory:  "files",
		Auth: AuthConfig{
			"HERE_MUST_BE_SECRET_KEY",
			time.Now().Add(1 * time.Hour).Unix(),
		},
	}
}
