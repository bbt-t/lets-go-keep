package main

import (
	"github.com/bbt-t/lets-go-keep/internal/app/client"
	"github.com/bbt-t/lets-go-keep/internal/config"
	"github.com/bbt-t/lets-go-keep/internal/controller/handlers"

	_ "github.com/jackc/pgx/v5/stdlib"
	log "github.com/sirupsen/logrus"
)

func main() {
	cfg := config.NewClientConfig()
	c := handlers.NewClientConnection(cfg.ServerAddress)
	h := handlers.NewClientHandlers(c)

	tui := client.NewTUI(h)

	log.Fatalln(tui.Run())
}
