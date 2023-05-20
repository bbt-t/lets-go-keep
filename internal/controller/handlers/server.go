package handlers

import (
	"context"
	"log"

	"github.com/bbt-t/lets-go-keep/internal/controller"
	"github.com/bbt-t/lets-go-keep/internal/entity"
	"github.com/bbt-t/lets-go-keep/internal/storage"
	"github.com/bbt-t/lets-go-keep/pkg"
)

// server struct for server handlers.
type server struct {
	Storage       storage.Storager
	Authenticator Authenticator
}

// newServerHandlers returns server handlers based on storage and authenticator (interface).
func newServerHandlers(s storage.Storager, a Authenticator) *server {
	return &server{
		Storage:       s,
		Authenticator: a,
	}
}

// LoginUser logins user by login and password.
func (s *server) LoginUser(credentials entity.UserCredentials) (entity.AuthToken, error) {
	if credentials.Login == "" || credentials.Password == "" {
		return "", controller.ErrFieldIsEmpty
	}

	credentials.Password = pkg.PasswordHash(credentials)

	userID, err := s.Storage.LoginUser(credentials)
	if err != nil {
		return "", err
	}

	authToken, errCreateToken := s.Authenticator.CreateToken(userID)
	if errCreateToken != nil {
		log.Println("Failed create authToken:", err)
		return "", storage.ErrUnknown
	}

	return authToken, nil
}

// CreateUser creates new user by login and password.
func (s *server) CreateUser(credentials entity.UserCredentials) (entity.AuthToken, error) {
	if credentials.Login == "" || credentials.Password == "" {
		return "", controller.ErrFieldIsEmpty
	}

	if err := s.Storage.CreateUser(entity.UserCredentials{
		Login:    credentials.Login,
		Password: pkg.PasswordHash(credentials),
	}); err != nil {
		return "", err
	}

	return s.LoginUser(credentials)
}

// GetRecordsInfo gets all records from storage.
func (s *server) GetRecordsInfo(ctx context.Context) ([]entity.Record, error) {
	userID, err := s.userValidate(ctx)
	if err != nil {
		return nil, err
	}

	return s.Storage.GetRecordsInfo(context.WithValue(ctx, "userID", userID))
}

// GetRecord get record from storage by ID.
func (s *server) GetRecord(ctx context.Context, recordID string) (entity.Record, error) {
	userID, err := s.userValidate(ctx)
	if err != nil {
		return entity.Record{}, err
	}

	return s.Storage.GetRecord(context.WithValue(ctx, "userID", userID), recordID)
}

// CreateRecord added record to storage.
func (s *server) CreateRecord(ctx context.Context, record entity.Record) error {
	userID, err := s.userValidate(ctx)
	if err != nil {
		return err
	}

	_, err = s.Storage.CreateRecord(context.WithValue(ctx, "userID", userID), record)
	return err
}

// DeleteRecord deletes record from storage.
func (s *server) DeleteRecord(ctx context.Context, recordID string) error {
	userID, err := s.userValidate(ctx)
	if err != nil {
		return err
	}

	return s.Storage.DeleteRecord(context.WithValue(ctx, "userID", userID), recordID)
}

// userValidate validate logic.
func (s *server) userValidate(ctx context.Context) (entity.UserID, error) {
	var userID entity.UserID

	token, ok := ctx.Value("authToken").(entity.AuthToken)
	if !ok {
		return userID, storage.ErrUnauthenticated
	}

	userIDValid, err := s.Authenticator.ValidateToken(token)
	if err != nil {
		return userID, err
	}
	return userIDValid, nil
}
