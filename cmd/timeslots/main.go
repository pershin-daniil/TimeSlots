package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pershin-daniil/TimeSlots/pkg/logger"
	"github.com/pershin-daniil/TimeSlots/pkg/telegram"
)

var tgToken = os.Getenv("TG_TOKEN")

func main() {
	log := logger.New()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tg, err := telegram.New(log, tgToken)
	if err != nil {
		log.Panic(err)
	}

	go func() {
		signCh := make(chan os.Signal, 1)
		signal.Notify(signCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
		<-signCh
		log.Infof("Received signal, shutting down...")
		cancel()
	}()

	if err = tg.Run(ctx); err != nil {
		log.Warn(err)
	}
}
