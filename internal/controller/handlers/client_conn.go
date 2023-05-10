package handlers

import (
	"context"
	"log"

	"github.com/bbt-t/lets-go-keep/internal/controller"
	"github.com/bbt-t/lets-go-keep/internal/entity"
	"github.com/bbt-t/lets-go-keep/internal/storage"
	pb "github.com/bbt-t/lets-go-keep/protocols/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// ClientConnGPRC keeps connection with server. Uses gRPC.
type ClientConnGPRC struct {
	pb.GophkeeperClient
}

// NewClientConnection connects to server and returning connection.
func newClientConn(serverAddress string) *ClientConnGPRC {
	conn, err := grpc.Dial(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	return &ClientConnGPRC{
		GophkeeperClient: pb.NewGophkeeperClient(conn),
	}
}

// Login logins user by login and password.
func (c *ClientConnGPRC) Login(credentials entity.UserCredentials) (string, error) {
	session, err := c.GophkeeperClient.Login(context.Background(), &pb.UserCredentials{
		Login:    credentials.Login,
		Password: credentials.Password,
	})

	code := status.Code(err)

	if code == codes.Unauthenticated {
		return "", storage.ErrWrongCredentials
	}
	if code == codes.Internal {
		return "", storage.ErrUnknown
	}
	if code == codes.InvalidArgument {
		return "", controller.ErrFieldIsEmpty
	}

	if err != nil {
		return "", err
	}

	return session.SessionToken, nil
}

// Register creates new user by login and password.
func (c *ClientConnGPRC) Register(credentials entity.UserCredentials) (string, error) {
	session, err := c.GophkeeperClient.Register(context.Background(), &pb.UserCredentials{
		Login:    credentials.Login,
		Password: credentials.Password,
	})

	code := status.Code(err)

	switch code {
	case codes.AlreadyExists:
		return "", storage.ErrLoginExists
	case codes.Internal:
		return "", storage.ErrUnknown
	case codes.InvalidArgument:
		return "", controller.ErrFieldIsEmpty
	}

	return session.SessionToken, nil
}

// GetRecordsInfo gets all record.
func (c *ClientConnGPRC) GetRecordsInfo(token entity.AuthToken) ([]entity.Record, error) {
	ctx := metadata.AppendToOutgoingContext(context.Background(), "authToken", string(token))
	gotRecords, err := c.GophkeeperClient.GetRecordsInfo(ctx, &emptypb.Empty{})
	code := status.Code(err)

	switch code {
	case codes.Internal:
		return nil, storage.ErrUnknown
	case codes.Unauthenticated:
		return nil, storage.ErrUserUnauthorized
	}

	records := make([]entity.Record, 0, len(gotRecords.Records))

	for _, record := range gotRecords.Records {
		records = append(records, entity.Record{
			ID:       record.Id,
			Metadata: record.Metadata,
			Type:     entity.RecordType(record.Type),
		})
	}

	return records, nil
}

// GetRecord gets record from server by ID.
func (c *ClientConnGPRC) GetRecord(token entity.AuthToken, recordID string) (entity.Record, error) {
	ctx := metadata.AppendToOutgoingContext(context.Background(), "authToken", string(token))
	gotRecord, err := c.GophkeeperClient.GetRecord(ctx, &pb.RecordID{
		Id: recordID,
	})
	record, code := entity.Record{}, status.Code(err)

	switch code {
	case codes.Internal:
		return record, storage.ErrUnknown
	case codes.Unauthenticated:
		return record, storage.ErrUserUnauthorized
	case codes.NotFound:
		return record, storage.ErrNotFound
	}

	record = entity.Record{
		ID:       gotRecord.Id,
		Metadata: gotRecord.Metadata,
		Type:     entity.RecordType(gotRecord.Type),
		Data:     gotRecord.StoredData,
	}
	return record, nil
}

// DeleteRecord deletes record from server by ID.
func (c *ClientConnGPRC) DeleteRecord(token entity.AuthToken, recordID string) error {
	ctx := metadata.AppendToOutgoingContext(context.Background(), "authToken", string(token))
	_, err := c.GophkeeperClient.DeleteRecord(ctx, &pb.RecordID{
		Id: recordID,
	})
	code := status.Code(err)

	switch code {
	case codes.Internal:
		return storage.ErrUnknown
	case codes.Unauthenticated:
		return storage.ErrUserUnauthorized
	case codes.NotFound:
		return storage.ErrNotFound
	}

	return nil
}

// CreateRecord creates record and saves to server.
func (c *ClientConnGPRC) CreateRecord(token entity.AuthToken, record entity.Record) error {
	ctx := metadata.AppendToOutgoingContext(context.Background(), "authToken", string(token))
	_, err := c.GophkeeperClient.CreateRecord(ctx, &pb.Record{
		Type:       pb.MessageType(record.Type),
		Metadata:   record.Metadata,
		StoredData: record.Data,
	})

	switch status.Code(err) {
	case codes.Internal:
		return storage.ErrUnknown
	case codes.Unauthenticated:
		return storage.ErrUserUnauthorized
	}

	return nil
}
