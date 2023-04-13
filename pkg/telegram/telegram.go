package telegram

import (
	"context"
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pershin-daniil/TimeSlots/pkg/models"
	"github.com/sirupsen/logrus"
)

type App interface {
	User(ctx context.Context, newUser models.UserRequest) (models.User, error)
	Status(ctx context.Context, userID int64) (string, error)
	CalendarApp
}

type CalendarApp interface {
	Events() []models.Event
}

type handlerFunc func(ctx context.Context, update tg.Update) error

type Telegram struct {
	log        *logrus.Entry
	app        App
	bot        *tg.BotAPI
	handlerMap map[string]handlerFunc
}

func New(log *logrus.Logger, app App, token string) (*Telegram, error) {
	bot, err := tg.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("failed to init bot: %w", err)
	}
	log.Debugf("Authorized on account %s", bot.Self.UserName)
	return &Telegram{
		log:        log.WithField("module", "telegram"),
		app:        app,
		bot:        bot,
		handlerMap: make(map[string]handlerFunc),
	}, nil
}

func (t *Telegram) Run(ctx context.Context) error {
	t.log.Infof("start listening for updates")
	u := tg.NewUpdate(0)
	u.Timeout = 60
	updates := t.bot.GetUpdatesChan(u)
	t.initHandlers()
	for {
		select {
		case <-ctx.Done():
			return nil
		case update := <-updates:
			t.processUpdate(ctx, update)
		}
	}
}

var (
	cmdStart     = "/start"
	cmdFT        = "/FT"
	cmdAbout     = "/about"
	cmdContact   = "/contact"
	cmdEducation = "/education"
	cmdPrice     = "/price"

	msgStartf = `
Привет, %s!

Это тренер Аня из DDX 👋 

Здесь можно посмотреть информацию обо мне, записаться на первую встречу и узнать стоимость занятий 🙂`

	msgAbout = `
Меня зовут Аня, и я - персональный тренер DDX Авиапарк ☺️

Когда-то давно я пришла в зал и... испугалась всех этих непонятных и одинаковых тренажёров. Я провела целый год на беговой дорожке, избегая их, но со временем мой страх перерос в желание узнать больше о спорте и тренировках.

Теперь я - тренер, которому нравится общаться с новыми людьми и помогать им достигать своих целей. Я знаю, как сложно начать свой путь к здоровому образу жизни, и я здесь, чтобы помочь вам.
Мой стиль - это персональный подход: всегда уделяю внимание индивидуальным потребностям каждого из моих подопечных 👍`

	msgContact = `
Ты можешь написать в телеграм 

Подписывайся на мой инстаграм`
	msgFT = `
OK 👌 Ждем подтверждения тренера.`

	msgPrice     = "50 рублей"
	msgEducation = "Тютюрский университет"

	btnBack = "Назад"

	btnAbout = "Обо мне"
	btnFT    = "Первая тренировка"

	btnContact   = "Контакты"
	btnEducation = "Образование"
	btnPrice     = "Прайс"

	btnPickSlot = "Выбрать"
)

var (
	welcomeKbd = tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(btnAbout, cmdAbout)),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(btnFT, cmdFT)))
	aboutKbd = tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(btnEducation, cmdEducation)),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(btnPrice, cmdPrice)),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(btnContact, cmdContact)),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(btnBack, cmdStart)))
	backToAboutKbd = tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(btnBack, cmdAbout)))
	contactKbd = tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonURL("Телеграм", "https://t.me/Filianan")),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonURL("Инстаграм", "https://www.instagram.com/filianan")),
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(btnBack, cmdAbout)))
	pickSlotKbd = tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData(btnPickSlot, "")))
)

func (t *Telegram) initHandlers() {
	t.handle(cmdStart, t.startHandler)
	t.handle(cmdAbout, t.aboutHandler)
	t.handle(cmdContact, t.contactHandler)
	t.handle(cmdPrice, t.priceHandler)
	t.handle(cmdEducation, t.educationHandler)
	t.handle(cmdFT, t.FTHandler)
}

func (t *Telegram) handle(command string, handler handlerFunc) {
	t.handlerMap[command] = handler
}

func (t *Telegram) processUpdate(ctx context.Context, update tg.Update) {
	var command string
	switch {
	case update.Message != nil && update.Message.IsCommand():
		currentMsg := tg.NewDeleteMessage(update.FromChat().ID, update.Message.MessageID)
		oldMsg := tg.NewDeleteMessage(update.FromChat().ID, update.Message.MessageID-1)
		t.bot.Send(currentMsg)
		t.bot.Send(oldMsg)
		command = update.Message.Text
	case update.Message != nil:
		currentMsg := tg.NewDeleteMessage(update.FromChat().ID, update.Message.MessageID)
		oldMsg := tg.NewDeleteMessage(update.FromChat().ID, update.Message.MessageID-1)
		t.bot.Send(currentMsg)
		t.bot.Send(oldMsg)
		command = cmdStart
	case update.CallbackQuery != nil:
		callback := tg.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
		if _, err := t.bot.Request(callback); err != nil {
			t.log.Warnf("err getting callback: %v", err)
		}
		command = update.CallbackQuery.Data
	}
	handler, ok := t.handlerMap[command]
	if ok {
		go func() {
			if err := handler(ctx, update); err != nil {
				t.log.Warnf("err during handle command: %v", err)
			}
		}()
	} else {
		t.log.Warnf("unknown command: %s, %#v", command, update)
	}
}

func (t *Telegram) parseUser(update tg.Update) models.UserRequest {
	return models.UserRequest{
		ID:        update.SentFrom().ID,
		LastName:  update.SentFrom().LastName,
		FirstName: update.SentFrom().FirstName,
		Status:    models.StatusUserGuest,
	}
}

func (t *Telegram) startHandler(ctx context.Context, update tg.Update) error {
	parsedUser := t.parseUser(update)

	user, err := t.app.User(ctx, parsedUser)
	if err != nil {
		return err
	}

	var msgText string
	switch user.Status {
	case models.StatusUserGuest:
		msgText = fmt.Sprintf(msgStartf, user.FirstName)
	case models.StatusUserFT:
		msgText = fmt.Sprintf(msgStartf+"FT", user.FirstName)
	}

	if update.CallbackQuery != nil {
		msgEdit := tg.NewEditMessageTextAndMarkup(update.FromChat().ID, update.CallbackQuery.Message.MessageID, msgText, welcomeKbd)
		if _, err = t.bot.Send(msgEdit); err != nil {
			return err
		}
		return nil
	}

	msg := tg.NewMessage(update.Message.Chat.ID, msgText)
	msg.ReplyMarkup = welcomeKbd
	if _, err = t.bot.Send(msg); err != nil {
		return err
	}
	return nil
}

func (t *Telegram) aboutHandler(ctx context.Context, update tg.Update) error {
	msg := tg.NewEditMessageTextAndMarkup(update.FromChat().ID, update.CallbackQuery.Message.MessageID, msgAbout, aboutKbd)
	if _, err := t.bot.Send(msg); err != nil {
		return err
	}
	return nil
}

func (t *Telegram) contactHandler(ctx context.Context, update tg.Update) error {
	msg := tg.NewEditMessageTextAndMarkup(update.FromChat().ID, update.CallbackQuery.Message.MessageID, msgContact, contactKbd)
	if _, err := t.bot.Send(msg); err != nil {
		return err
	}
	return nil
}

func (t *Telegram) priceHandler(ctx context.Context, update tg.Update) error {
	msg := tg.NewEditMessageTextAndMarkup(update.FromChat().ID, update.CallbackQuery.Message.MessageID, msgPrice, backToAboutKbd)
	if _, err := t.bot.Send(msg); err != nil {
		return err
	}
	return nil
}

func (t *Telegram) educationHandler(ctx context.Context, update tg.Update) error {
	msg := tg.NewEditMessageTextAndMarkup(update.FromChat().ID, update.CallbackQuery.Message.MessageID, msgEducation, backToAboutKbd)
	if _, err := t.bot.Send(msg); err != nil {
		return err
	}
	return nil
}

func (t *Telegram) FTHandler(ctx context.Context, update tg.Update) error {
	msg := tg.NewMessage(update.FromChat().ID, "Доступное время для записи на первую тренировку 👇")
	events := t.app.Events()
	_, _ = t.bot.Send(msg)
	for i := 0; i < 1; i++ {
		event := events[i]
		msg = tg.NewMessage(update.FromChat().ID, fmt.Sprintf("Название: %s\nНачало: %s\nКонец: %s", event.Title, event.Start, event.End))
		msg.ReplyMarkup = pickSlotKbd
		_, _ = t.bot.Send(msg)
	}

	callback := tg.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
	if _, err := t.bot.Request(callback); err != nil {
		panic(err)
	}

	msg = tg.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
	if _, err := t.bot.Send(msg); err != nil {
		panic(err)

	}
	return nil
}
