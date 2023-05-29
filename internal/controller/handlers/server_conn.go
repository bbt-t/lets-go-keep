package handlers

import (
	"context"
	"errors"
	"net"

	"github.com/bbt-t/lets-go-keep/internal/controller"
	"github.com/bbt-t/lets-go-keep/internal/entity"
	"github.com/bbt-t/lets-go-keep/internal/storage"
	pb "github.com/bbt-t/lets-go-keep/protocols/grpc"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// ServerConn keeps server endpoints alive.
type ServerConn struct {
	pb.UnimplementedGophkeeperServer
	Handlers ServerHandlers
	server   *grpc.Server
}

// NewServerConn returns new server connection.
func NewServerConn(h ServerHandlers) *ServerConn {
	return &ServerConn{
		Handlers: h,
	}
}

// Run runs server listener.
func (s *ServerConn) Run(_ context.Context, runAddress string) {
	listen, err := net.Listen("tcp", runAddress)
	if err != nil {
		log.Fatal(err)
	}

	grpcServ := grpc.NewServer()
	pb.RegisterGophkeeperServer(grpcServ, s)

	go func() {
		log.Println("Сервер gRPC начал работу")
		if err := grpcServ.Serve(listen); err != nil {
			log.Fatal(err)
		}
	}()

	s.server = grpcServ
}

func (s *ServerConn) Stop() {
	s.server.GracefulStop()
	log.Println("Shutdown server gracefully.")
}

// Register process register endpoint.
func (s *ServerConn) Register(_ context.Context, credentials *pb.UserCredentials) (*pb.Session, error) {
	token, err := s.Handlers.CreateUser(entity.UserCredentials{
		Login:    credentials.Login,
		Password: credentials.Password,
	})

	if errors.Is(err, controller.ErrFieldIsEmpty) {
		log.Infoln(err)

		return nil, status.Errorf(codes.InvalidArgument, "Login or password is empty.")
	}

	if errors.Is(err, storage.ErrLoginExists) {
		log.Infoln(err)

		return nil, status.Errorf(codes.AlreadyExists, "Login already exists.")
	}

	if err != nil {
		log.Warnf("%s %s :: %v", "register new user fault", credentials.Login, err)

		return nil, status.Errorf(codes.Internal, "Internal server error.")
	}

	return &pb.Session{SessionToken: string(token)}, nil
}

// Login process login endpoint.
func (s *ServerConn) Login(_ context.Context, credentials *pb.UserCredentials) (*pb.Session, error) {
	token, err := s.Handlers.LoginUser(entity.UserCredentials{
		Login:    credentials.Login,
		Password: credentials.Password,
	})

	if errors.Is(err, controller.ErrFieldIsEmpty) {
		log.Infoln(err)

		return nil, status.Errorf(codes.InvalidArgument, "Login or password is empty.")
	}

	if errors.Is(err, storage.ErrWrongCredentials) {
		log.Infoln(err)

		return nil, status.Errorf(codes.Unauthenticated, "Wrong login or password.")
	}

	if err != nil {
		log.Warnf("%s %s :: %v", "login fault", credentials.Login, err)

		return nil, status.Errorf(codes.Internal, "Internal server error.")
	}

	return &pb.Session{SessionToken: string(token)}, nil
}

// GetRecordsInfo process get all records endpoint.
func (s *ServerConn) GetRecordsInfo(ctx context.Context, _ *emptypb.Empty) (*pb.RecordsList, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok || len(md.Get("authToken")) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "Didn't send metadata for authentication.")
	}

	token := entity.AuthToken(md.Get("authToken")[0])
	ctx = context.WithValue(ctx, "authToken", token)

	records, err := s.Handlers.GetRecordsInfo(ctx)

	if errors.Is(err, storage.ErrUnauthenticated) {
		log.Infoln(err)

		return nil, status.Errorf(codes.Unauthenticated, "Bad authentication token.")
	}

	if err != nil {
		log.Warnf("%s :: %v", "get record info fault", err)

		return nil, status.Errorf(codes.Internal, "Internal server error.")
	}

	recordsList := make([]*pb.Record, 0, len(records))

	for _, record := range records {
		recordsList = append(recordsList, &pb.Record{
			Id:       record.ID,
			Metadata: record.Metadata,
			Type:     pb.MessageType(record.Type),
		})
	}

	return &pb.RecordsList{Records: recordsList}, nil
}

// GetRecord process get record endpoint.
func (s *ServerConn) GetRecord(ctx context.Context, recordID *pb.RecordID) (*pb.Record, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok || len(md.Get("authToken")) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "Didn't send metadata for authentication.")
	}

	token := entity.AuthToken(md.Get("authToken")[0])
	ctx = context.WithValue(ctx, "authToken", token)

	record, err := s.Handlers.GetRecord(ctx, recordID.Id)

	if errors.Is(err, storage.ErrUnauthenticated) {
		log.Infoln(err)

		return nil, status.Errorf(codes.Unauthenticated, "Bad authentication token.")
	}

	if errors.Is(err, storage.ErrNotFound) {
		log.Infoln(err)

		return nil, status.Errorf(codes.NotFound, "Not found record with such id.")
	}

	if err != nil {
		log.Warnf("%s :: %v", "get record fault", err)

		return nil, status.Errorf(codes.Internal, "Internal server error.")
	}

	return &pb.Record{
		Id:         record.ID,
		Type:       pb.MessageType(record.Type),
		Metadata:   record.Metadata,
		StoredData: record.Data,
	}, nil
}

// CreateRecord process create record endpoint.
func (s *ServerConn) CreateRecord(ctx context.Context, record *pb.Record) (*emptypb.Empty, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok || len(md.Get("authToken")) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "Didn't send metadata for authentication.")
	}

	token := entity.AuthToken(md.Get("authToken")[0])
	ctx = context.WithValue(ctx, "authToken", token)

	err := s.Handlers.CreateRecord(ctx, entity.Record{
		Metadata: record.Metadata,
		Type:     entity.RecordType(record.Type),
		Data:     record.StoredData,
	})

	if errors.Is(err, storage.ErrUnauthenticated) {
		log.Infoln(err)

		return &emptypb.Empty{}, status.Errorf(codes.Unauthenticated, "Bad authentication token.")
	}

	if err != nil {
		log.Warnf("%s :: %v", "create record fault", err)

		return &emptypb.Empty{}, status.Errorf(codes.Internal, "Internal server error.")
	}

	return &emptypb.Empty{}, nil
}

// DeleteRecord process delete record endpoint.
func (s *ServerConn) DeleteRecord(ctx context.Context, recordID *pb.RecordID) (*emptypb.Empty, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok || len(md.Get("authToken")) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "Didn't send metadata for authentication.")
	}

	token := entity.AuthToken(md.Get("authToken")[0])
	ctx = context.WithValue(ctx, "authToken", token)

	err := s.Handlers.DeleteRecord(ctx, recordID.Id)

	if errors.Is(err, storage.ErrUnauthenticated) {
		log.Infoln(err)

		return &emptypb.Empty{}, status.Errorf(codes.Unauthenticated, "Bad authentication token.")
	}

	if errors.Is(err, storage.ErrNotFound) {
		log.Infoln(err)

		return nil, status.Errorf(codes.NotFound, "Not found record with such id.")
	}

	if err != nil {
		log.Warnf("%s :: %v", "delete record fault", err)

		return &emptypb.Empty{}, status.Errorf(codes.Internal, "Internal server error.")
	}

	return &emptypb.Empty{}, nil
}
