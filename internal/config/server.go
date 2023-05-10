package config

import "time"

// ServerConfig struct for server config.
type ServerConfig struct {
	RunAddress      string
	DBConnectionURL string
	FilesDirectory  string
	Auth            AuthConfig
}

type AuthConfig struct {
	SecretJWT      []byte
	ExpirationTime int64
}

// NewServerConfig gets server config.
func NewServerConfig() ServerConfig {
	auth := AuthConfig{
		SecretJWT:      []byte("HERE_MUST_BE_SECRET_KEY"),
		ExpirationTime: time.Now().Add(1 * time.Hour).Unix(),
	}
	return ServerConfig{
		RunAddress:      ":3200",
		DBConnectionURL: "",
		FilesDirectory:  "files",
		Auth:            auth,
	}
}
