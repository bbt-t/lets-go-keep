package storage

import (
	"context"
	"errors"
	"io"
	"os"

	"github.com/bbt-t/lets-go-keep/internal/entity"

	log "github.com/sirupsen/logrus"
)

// fileStorage keeps records on disk.
type fileStorage struct {
	directory string
}

// newFileStorage returns new file storage.
func newFileStorage(directory string) *fileStorage {
	if err := os.Mkdir(directory, os.ModePerm); err != nil && !os.IsExist(err) {
		log.Fatalln("Failed open directory for file storage")

		return nil
	}

	return &fileStorage{directory: directory}
}

// GetRecord reads file with record data.
func (storage *fileStorage) GetRecord(ctx context.Context, recordID string) (entity.Record, error) {
	metadata, ok := ctx.Value("recordMetadata").(string)
	if !ok {
		log.Println("Failed get record metadata from context in getting file record")
		return entity.Record{}, ErrUnknown
	}

	file, err := os.Open(storage.directory + "/" + recordID)
	if errors.Is(err, os.ErrNotExist) {
		log.Infoln(err)

		return entity.Record{}, ErrNotFound
	}
	if err != nil {
		log.Infoln(err)

		return entity.Record{}, ErrUnknown
	}

	data, errReadAll := io.ReadAll(file)
	if errReadAll != nil {
		log.Infoln(err)

		return entity.Record{}, ErrUnknown
	}

	return entity.Record{
		ID:       recordID,
		Metadata: metadata,
		Type:     entity.TypeFile,
		Data:     data,
	}, nil
}

// DeleteRecord deletes file with record data.
func (storage *fileStorage) DeleteRecord(_ context.Context, recordID string) error {
	filename := storage.directory + "/" + recordID
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		return ErrNotFound
	}

	if err := os.RemoveAll(filename); err != nil {
		log.Infoln(err)

		return ErrUnknown
	}

	return nil
}

// CreateRecord creates new file with record data.
func (storage *fileStorage) CreateRecord(_ context.Context, record entity.Record) (string, error) {
	file, err := os.Create(storage.directory + "/" + record.ID)
	if err != nil {
		log.Infoln(err)

		return "", ErrUnknown
	}

	if _, err := file.Write(record.Data); err != nil {
		return "", ErrUnknown
	}

	return record.ID, nil
}
