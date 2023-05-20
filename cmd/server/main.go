package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/bbt-t/lets-go-keep/internal/config"
	"github.com/bbt-t/lets-go-keep/internal/controller/handlers"
	"github.com/bbt-t/lets-go-keep/internal/storage"

	_ "github.com/jackc/pgx/v5/stdlib"
	log "github.com/sirupsen/logrus"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.WarnLevel)
}

func main() {
	cfg := config.NewServerConfig()

	db := storage.NewDBStorage(cfg.DBConnectionURL)
	db.MigrateUP()

	files := storage.NewFileStorage(cfg.FilesDirectory)
	s := storage.NewStorage(db, files)

	jwtAuth := handlers.NewAuthenticatorJWT([]byte(cfg.Auth.SecretJWT), cfg.Auth.ExpirationTime)
	h := handlers.NewServerHandlers(s, jwtAuth)
	server := handlers.NewServerConn(h)

	go server.Run(context.Background(), cfg.RunAddress)

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-sigint

	server.Stop()
}
