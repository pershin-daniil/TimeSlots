package telegram

import (
	"context"
	"fmt"
	"time"

	"github.com/pershin-daniil/TimeSlots/pkg/calendar"

	"github.com/pershin-daniil/TimeSlots/pkg/models"

	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

type Telegram struct {
	log *logrus.Entry
	bot *tele.Bot
	cal Calendar
}

type Calendar interface {
	Events() []models.Event
}

func New(log *logrus.Logger, token string, cal *calendar.Calendar) (*Telegram, error) {
	config := tele.Settings{
		Token:     token,
		Poller:    &tele.LongPoller{Timeout: 10 * time.Second},
		ParseMode: tele.ModeMarkdown,
		OnError: func(err error, ctx tele.Context) {
			log.Panicf("catch err: %v", err)
		},
	}

	b, err := tele.NewBot(config)
	if err != nil {
		return nil, fmt.Errorf("new bot faild: %w", err)
	}

	t := Telegram{
		log: log.WithField("module", "telegram"),
		bot: b,
		cal: cal,
	}

	t.initButtons()
	t.initHandlers()
	return &t, nil
}

func (t *Telegram) Run(ctx context.Context) {
	go func() {
		<-ctx.Done()
		t.bot.Stop()
	}()
	t.log.Infof("Starting telegram bot as %v", t.bot.Me.Username)
	t.bot.Start()
}
