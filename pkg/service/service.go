package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/pershin-daniil/TimeSlots/pkg/models"
	"github.com/sirupsen/logrus"
)

type Calendar interface {
	SlotsOneWeek() []models.Event
}

type Store interface {
	User(ctx context.Context, userID int64) (models.User, error)
	Status(ctx context.Context, userID int64) (string, error)
	CreateUser(ctx context.Context, user models.UserRequest) (models.User, error)
	UpdateUser(ctx context.Context, user models.UserRequest) (models.User, error)
}

type App struct {
	log   *logrus.Entry
	cal   Calendar
	store Store
}

func New(log *logrus.Logger, cal Calendar, store Store) *App {
	return &App{
		log:   log.WithField("module", "service"),
		cal:   cal,
		store: store,
	}
}

func (a *App) User(ctx context.Context, newUser models.UserRequest) (models.User, error) {
	user, err := a.store.User(ctx, newUser.ID)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		user, err = a.store.CreateUser(ctx, newUser)
	case err != nil:
		return models.User{}, fmt.Errorf("service: %w", err)
	}

	var update bool
	switch {
	case user.LastName != newUser.LastName:
		update = true
	case user.FirstName != newUser.FirstName:
		update = true
	case user.Status != newUser.Status:
		update = true
	}

	if update {
		user, err = a.store.UpdateUser(ctx, newUser)
		if err != nil {
			return models.User{}, fmt.Errorf("service: %w", err)
		}
	}

	return user, nil
}

func (a *App) Status(ctx context.Context, userID int64) (string, error) {
	status, err := a.store.Status(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("service: %w", err)
	}
	return status, nil
}

func (a *App) Events() []models.Event {
	slots := a.cal.SlotsOneWeek()
	return slots
}
