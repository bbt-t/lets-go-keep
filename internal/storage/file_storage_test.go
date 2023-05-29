package storage

import (
	"context"
	"os"
	"testing"

	"github.com/bbt-t/lets-go-keep/internal/config"
	"github.com/bbt-t/lets-go-keep/internal/entity"

	"github.com/stretchr/testify/assert"
)

var filesDirectory = config.NewServerConfig().FilesDirectory

func TestNewFileStorage(t *testing.T) {
	assert.NotPanics(t, func() {
		newFileStorage(filesDirectory)
	})
}

func TestFileStorage_CreateRecord(t *testing.T) {
	storage := newFileStorage(filesDirectory)

	tc := []struct {
		name    string
		prepare func()
		valid   func()
	}{
		{
			"Create file record",
			func() {
				id, err := storage.CreateRecord(context.Background(), entity.Record{
					ID:   "1",
					Type: entity.TypeFile,
					Data: []byte("text"),
				})
				assert.NoError(t, err)
				assert.Equal(t, "1", id)
			},
			func() {
				assert.DirExists(t, filesDirectory)
				assert.FileExists(t, filesDirectory+"/1")
			},
		},
	}

	for _, test := range tc {
		t.Log(test.name)
		test.prepare()
		test.valid()
	}

	assert.NoError(t, os.RemoveAll(filesDirectory))
}

func TestFileStorage_GetRecord(t *testing.T) {
	storage := newFileStorage(filesDirectory)

	tc := []struct {
		name    string
		prepare func()
		valid   func()
	}{
		{
			"Get existed file record",
			func() {
				id, err := storage.CreateRecord(context.Background(), entity.Record{
					ID:   "1",
					Type: entity.TypeFile,
					Data: []byte("text"),
				})
				assert.NoError(t, err)
				assert.Equal(t, "1", id)
			},
			func() {
				ctx := context.WithValue(context.Background(), "recordMetadata", "file.txt")
				record, err := storage.GetRecord(ctx, "1")
				assert.NoError(t, err)

				assert.Equal(t, entity.Record{
					ID:       "1",
					Metadata: "file.txt",
					Type:     entity.TypeFile,
					Data:     []byte("text"),
				}, record)

			},
		},
		{
			"Get existed file record, but don't provide record metadata",
			func() {
				id, err := storage.CreateRecord(context.Background(), entity.Record{
					ID:   "1",
					Type: entity.TypeFile,
					Data: []byte("text"),
				})
				assert.NoError(t, err)
				assert.Equal(t, "1", id)
			},
			func() {
				record, err := storage.GetRecord(context.Background(), "1")
				assert.Equal(t, ErrUnknown, err)
				assert.Empty(t, record)
			},
		},
		{
			"Get non existed file record",
			func() {
				id, err := storage.CreateRecord(context.Background(), entity.Record{
					ID:   "1",
					Type: entity.TypeFile,
					Data: []byte("text"),
				})
				assert.NoError(t, err)
				assert.Equal(t, "1", id)
			},
			func() {
				ctx := context.WithValue(context.Background(), "recordMetadata", "file.txt")
				record, err := storage.GetRecord(ctx, "2")
				assert.Equal(t, ErrNotFound, err)
				assert.Empty(t, record)
			},
		},
	}

	for _, test := range tc {
		t.Log(test.name)
		test.prepare()
		test.valid()
	}

	assert.NoError(t, os.RemoveAll(filesDirectory))
}

func TestFileStorage_DeleteRecord(t *testing.T) {
	storage := newFileStorage(filesDirectory)

	tc := []struct {
		name    string
		prepare func()
		valid   func()
	}{
		{
			"Delete existed file record",
			func() {
				id, err := storage.CreateRecord(context.Background(), entity.Record{
					ID:   "1",
					Type: entity.TypeFile,
					Data: []byte("text"),
				})
				assert.NoError(t, err)
				assert.Equal(t, "1", id)
			},
			func() {
				err := storage.DeleteRecord(context.Background(), "1")
				assert.NoError(t, err)
				assert.NoFileExists(t, filesDirectory+"/1")
			},
		},
		{
			"Delete non existed file record",
			func() {
				id, err := storage.CreateRecord(context.Background(), entity.Record{
					ID:   "1",
					Type: entity.TypeFile,
					Data: []byte("text"),
				})
				assert.NoError(t, err)
				assert.Equal(t, "1", id)
			},
			func() {
				err := storage.DeleteRecord(context.Background(), "2")
				assert.Equal(t, ErrNotFound, err)
				assert.NoFileExists(t, filesDirectory+"/2")
			},
		},
	}

	for _, test := range tc {
		t.Log(test.name)
		test.prepare()
		test.valid()
	}

	assert.NoError(t, os.RemoveAll(filesDirectory))
}
