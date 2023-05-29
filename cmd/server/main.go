package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bbt-t/lets-go-keep/internal/config"
	"github.com/bbt-t/lets-go-keep/internal/controller/handlers"
	"github.com/bbt-t/lets-go-keep/internal/storage"

	_ "github.com/jackc/pgx/v5/stdlib"
	log "github.com/sirupsen/logrus"
)

func init() {
	fileName := time.Now().Format("01-02-2006")
	f, err := os.OpenFile(
		fmt.Sprintf("./logs/%s%s", fileName, ".log"),
		os.O_APPEND|os.O_CREATE|os.O_RDWR,
		0666,
	)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(f)
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
