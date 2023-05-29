package handlers

import (
	"context"

	"github.com/bbt-t/lets-go-keep/internal/storage"

	"github.com/bbt-t/lets-go-keep/internal/entity"
)

// ClientHandlers interface for Client.
type ClientHandlers interface {
	Login(credentials entity.UserCredentials) error
	Register(credentials entity.UserCredentials) error
	GetRecordsInfo() ([]entity.Record, error)
	GetRecord(recordID string) (entity.Record, error)
	CreateRecord(record entity.Record) error
	DeleteRecord(recordID string) error
}

// NewClientHandlers returns new client handlers (interface).
func NewClientHandlers(conn ClientConnection) ClientHandlers {
	return newClientHandlers(conn)
}

// Authenticator is interface for user authenticating. Should can creates tokens, and gets userIDs from them.
//
//go:generate mockery --name Authenticator
type Authenticator interface {
	CreateToken(userID entity.UserID) (entity.AuthToken, error)
	ValidateToken(token entity.AuthToken) (entity.UserID, error)
}

// NewAuthenticatorJWT gets new authenticatorJWT (interface).
func NewAuthenticatorJWT(secretKey []byte, expirationTime int64) Authenticator {
	return newAuthenticatorJWT(secretKey, expirationTime)
}

// ClientConnection describes client connection.
//
//go:generate mockery --name ClientConn
type ClientConnection interface {
	Login(credentials entity.UserCredentials) (string, error)
	Register(credentials entity.UserCredentials) (string, error)
	GetRecordsInfo(token entity.AuthToken) ([]entity.Record, error)
	GetRecord(token entity.AuthToken, recordID string) (entity.Record, error)
	DeleteRecord(token entity.AuthToken, recordID string) error
	CreateRecord(token entity.AuthToken, record entity.Record) error
}

// NewClientConnection connects to server and returning connection (interface).
func NewClientConnection(serverAddress string) ClientConnection {
	return newClientConn(serverAddress)
}

// ServerHandlers interface for server handlers
//
//go:generate mockery --name ServerHandlers
type ServerHandlers interface {
	LoginUser(credentials entity.UserCredentials) (entity.AuthToken, error)
	CreateUser(credentials entity.UserCredentials) (entity.AuthToken, error)
	GetRecordsInfo(ctx context.Context) ([]entity.Record, error)
	GetRecord(ctx context.Context, recordID string) (entity.Record, error)
	CreateRecord(ctx context.Context, record entity.Record) error
	DeleteRecord(ctx context.Context, recordID string) error
}

// NewServerHandlers returns server handlers based on storage and authenticator.
func NewServerHandlers(s storage.Storager, a Authenticator) ServerHandlers {
	return newServerHandlers(s, a)
}
