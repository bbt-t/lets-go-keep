package config

import (
	"flag"
	"log"
	"time"

	"github.com/caarlos0/env/v8"
)

// ServerConfig struct for server config.
type ServerConfig struct {
	RunAddress      string `env:"SERVER_PORT" envDefault:":3200"`
	DBConnectionURL string `env:"DATABASE_DSN"`
	FilesDirectory  string `env:"FILE_STORAGE_PATH" envDefault:"files"`
	Auth            AuthConfig
}

type AuthConfig struct {
	SecretJWT      []byte `env:"SECRET_JWT" envDefault:"HERE_MUST_BE_SECRET_KEY"`
	ExpirationTime int64
}

// NewServerConfig gets server config.
func NewServerConfig() ServerConfig {
	var (
		cfg       ServerConfig
		secretJWT string
	)

	flag.StringVar(&cfg.RunAddress, "a", "", "server address")
	flag.StringVar(&cfg.FilesDirectory, "f", "", "file db path")
	flag.StringVar(&cfg.DBConnectionURL, "d", "", "postgres DSN (url)")
	flag.StringVar(&secretJWT, "s", "", "secret key for JWT")

	if err := env.Parse(&cfg); err != nil {
		log.Printf("%+v\n", err)
	}
	if err := env.Parse(&secretJWT); err != nil {
		log.Printf("%+v\n", err)
	}
	flag.Parse()

	cfg.Auth = AuthConfig{
		SecretJWT:      []byte(secretJWT),
		ExpirationTime: time.Now().Add(1 * time.Hour).Unix(),
	}

	return cfg
}
