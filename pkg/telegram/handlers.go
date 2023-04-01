package telegram

import (
	"fmt"
	"time"

	"github.com/pershin-daniil/TimeSlots/pkg/models"

	tele "gopkg.in/telebot.v3"
)

var coachCode = "TRAIN"

func (t *Telegram) initHandlers() {
	t.bot.Handle(cmdStart, t.startHandler)
	t.bot.Handle(&registrationBtn, t.registrationHandler)
	t.bot.Handle(&availableMeetingsBtn, t.scheduleHandler)
	t.bot.Handle(&myMeetingBtn, t.meetingsHandler)
	t.bot.Handle(&notificationBtn, t.notifyHandler)
	t.bot.Handle(&cancelMeetingBtn, t.cancelMeetingHandler)
	t.bot.Handle(&pickSlotBtn, t.pickSlotHandler)
	t.bot.Handle(tele.OnText, t.textHandler)
}

func (t *Telegram) textHandler(ctx tele.Context) error {
	if ctx.Text() == coachCode {
		_ = ctx.Delete()
		_ = ctx.Send("У тебя самый лучший тренер, го посмотрим его расписание", availMeetings)
	} else {
		time.Sleep(5 * time.Second)
		_ = ctx.Delete()
	}
	return nil
}

func (t *Telegram) startHandler(ctx tele.Context) error {
	_ = ctx.Delete()
	user := parseUser(ctx)
	msg := `Вступительная речь
Предложение продолжить ` + fmt.Sprintf("%v", user.ID)

	return ctx.Send(msg, registration)
}

func parseUser(ctx tele.Context) *models.User {
	return &models.User{
		ID:        ctx.Sender().ID,
		LastName:  ctx.Sender().LastName,
		FirstName: ctx.Sender().FirstName,
	}
}

func (t *Telegram) registrationHandler(ctx tele.Context) error {
	msg := "Введите код тренера"
	return ctx.Edit(msg)
}

func (t *Telegram) scheduleHandler(ctx tele.Context) error {
	_, _ = t.bot.Send(ctx.Sender(), "Здесь доступное расписание тренера, на которое можно записаться.")
	events := t.cal.Events()
	t.log.Infof("%s", events)
	if len(events) == 0 {
		return ctx.Send("Свободных слотов нет", showMeetings)
	}
	for _, event := range events {
		_, _ = t.bot.Send(ctx.Sender(), fmt.Sprintf("Name: %s\nStart: %s\nEnd: %s", event.Title, event.Start, event.End), pickSlot)
	}
	return ctx.Send("Когда закончишь выбирать, нажимай сюда 😄", showMeetings)
}

func (t *Telegram) pickSlotHandler(ctx tele.Context) error {
	msg := `OK👌. Ждем подтверждения тренера.`

	defer func() {
		time.Sleep(5 * time.Second)
		_ = ctx.Delete()
	}()

	return ctx.Edit(msg)
}

func (t *Telegram) meetingsHandler(ctx tele.Context) error {
	// TODO: отображение расписания
	msg := "Моё расписание"
	return ctx.Edit(msg, settings)
}

func (t *Telegram) notifyHandler(ctx tele.Context) error {
	// TODO: notify logic
	msg := "Настройка нотификации"
	return ctx.Edit(msg, showMeetings)
}

func (t *Telegram) cancelMeetingHandler(ctx tele.Context) error {
	// TODO: canceling meeting
	msg := "Отмена занятия"
	return ctx.Edit(msg, showMeetings)
}
