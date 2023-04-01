package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/pershin-daniil/TimeSlots/pkg/calendar"
	"github.com/pershin-daniil/TimeSlots/pkg/logger"
	"github.com/pershin-daniil/TimeSlots/pkg/telegram"
)

var tgToken = os.Getenv("TG_TOKEN")

func main() {
	log := logger.New()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cal := calendar.New(ctx, log)
	tg, err := telegram.New(log, tgToken, cal)
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
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		tg.Run(ctx)
	}()
	wg.Wait()
}
