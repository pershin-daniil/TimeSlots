package main

import (
	"context"
	"github.com/pershin-daniil/TimeSlots/pkg/pgstore"
	migrate "github.com/rubenv/sql-migrate"
	"os"
	"os/signal"
	"syscall"

	"github.com/pershin-daniil/TimeSlots/pkg/service"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pershin-daniil/TimeSlots/pkg/calendar"
	"github.com/pershin-daniil/TimeSlots/pkg/logger"
	"github.com/pershin-daniil/TimeSlots/pkg/telegram"
)

var (
	tgToken = os.Getenv("TG_TOKEN")
	dsn     = os.Getenv("PG_DSN")
)

func main() {
	log := logger.New()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cal := calendar.New(ctx, log)
	store, err := pgstore.New(ctx, log, dsn)
	if err != nil {
		log.Panic(err)
	}
	if err = store.Migrate(migrate.Up); err != nil {
		log.Panic(err)
	}
	app := service.New(log, cal, store)
	tg, err := telegram.New(log, app, tgToken)
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
