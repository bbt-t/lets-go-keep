package storage

import (
	"context"

	"github.com/bbt-t/lets-go-keep/internal/entity"
)

type DataBaseStorage interface {
	MigrateUP()
	CreateUser(credentials entity.UserCredentials) error
	LoginUser(credentials entity.UserCredentials) (entity.UserID, error)
	GetRecordsInfo(ctx context.Context) ([]entity.Record, error)
	CreateRecord(ctx context.Context, record entity.Record) (string, error)
	GetRecord(ctx context.Context, recordID string) (entity.Record, error)
	DeleteRecord(ctx context.Context, recordID string) error
}

// NewDBStorage connects to DB (interface).
func NewDBStorage(connectionURL string) DataBaseStorage {
	return newDBStorage(connectionURL)
}

// FileStorager interface for storage, which can storage files.
//
//go:generate mockery --name FileStorager
type FileStorager interface {
	GetRecord(ctx context.Context, recordID string) (entity.Record, error)
	CreateRecord(ctx context.Context, record entity.Record) (string, error)
	DeleteRecord(ctx context.Context, recordID string) error
}

// NewFileStorage returns new file storage (interface).
func NewFileStorage(directory string) FileStorager {
	return newFileStorage(directory)
}

// Storager interface for storage, which can storage only text data.
//
//go:generate mockery --name Storager
type Storager interface {
	CreateUser(credentials entity.UserCredentials) error
	LoginUser(credentials entity.UserCredentials) (entity.UserID, error)
	GetRecordsInfo(ctx context.Context) ([]entity.Record, error)
	FileStorager
}
