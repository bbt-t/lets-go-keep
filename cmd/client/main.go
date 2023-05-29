package main

import (
	"fmt"
	"github.com/bbt-t/lets-go-keep/internal/app/client"
	"github.com/bbt-t/lets-go-keep/internal/config"
	"github.com/bbt-t/lets-go-keep/internal/controller/handlers"
	"os"
	"time"

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
	cfg := config.NewClientConfig()
	c := handlers.NewClientConnection(cfg.ServerAddress)
	h := handlers.NewClientHandlers(c)

	tui := client.NewTUI(h)

	log.Fatalln(tui.Run())
}
