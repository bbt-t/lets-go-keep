package storage

import (
	"context"
	"errors"

	"github.com/bbt-t/lets-go-keep/internal/entity"
)

// Storage struct which saves to DB and file storage.
type Storage struct {
	DBStorage   Storager
	FileStorage FileStorager
}

// NewStorage returns new storage.
func NewStorage(DBStorage Storager, fileStorage FileStorager) *Storage {
	return &Storage{
		DBStorage:   DBStorage,
		FileStorage: fileStorage,
	}
}

// CreateUser creates new user and saves to DB storage.
func (s *Storage) CreateUser(credentials entity.UserCredentials) error {
	return s.DBStorage.CreateUser(credentials)
}

// LoginUser check user login using DB storage.
func (s *Storage) LoginUser(credentials entity.UserCredentials) (entity.UserID, error) {
	return s.DBStorage.LoginUser(credentials)
}

// GetRecordsInfo gets all records from user from DB storage.
func (s *Storage) GetRecordsInfo(ctx context.Context) ([]entity.Record, error) {
	return s.DBStorage.GetRecordsInfo(ctx)
}

// CreateRecord creates record, saves to DB. If record type is file, saves to file storage too.
func (s *Storage) CreateRecord(ctx context.Context, record entity.Record) (string, error) {
	data := record.Data

	if record.Type == entity.TypeFile {
		record.Data = nil
	}

	id, err := s.DBStorage.CreateRecord(ctx, record)
	if err != nil {
		return "", err
	}

	if record.Type == entity.TypeFile {
		record.ID = id
		record.Data = data
		_, err = s.FileStorage.CreateRecord(ctx, record)
		return "", err
	}

	return id, nil
}

// DeleteRecord deletes record from DB storage. If record type is file, deletes from file storage too.
func (s *Storage) DeleteRecord(ctx context.Context, recordID string) error {
	err := s.DBStorage.DeleteRecord(ctx, recordID)
	if err != nil {
		return err
	}

	err = s.FileStorage.DeleteRecord(ctx, recordID)
	if !errors.Is(err, ErrNotFound) && err != nil {
		return ErrUnknown
	}

	return nil
}

// GetRecord gets record from DB or file storage.
func (s *Storage) GetRecord(ctx context.Context, recordID string) (entity.Record, error) {
	record, err := s.DBStorage.GetRecord(ctx, recordID)
	if err != nil {
		return record, err
	}

	if record.Type == entity.TypeFile {
		ctx = context.WithValue(ctx, "recordMetadata", record.Metadata)
		return s.FileStorage.GetRecord(ctx, recordID)
	}

	return record, nil
}
