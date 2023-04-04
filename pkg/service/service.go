package service

import (
	"context"
	"fmt"
	"github.com/pershin-daniil/TimeSlots/pkg/models"
	"github.com/sirupsen/logrus"
)

type Calendar interface {
	Events() []models.Event
}

type Store interface {
	User(ctx context.Context, user models.UserRequest) (models.User, error)
	Session(ctx context.Context, userID int64, state string) (models.Session, error)
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

func (s *App) User(newUser models.UserRequest) (models.User, error) {
	ctx := context.Background()
	user, err := s.store.User(ctx, newUser)
	if err != nil {
		return models.User{}, fmt.Errorf("service: %w", err)
	}
	return user, nil
}

func (s *App) Session(userID int64, state string) (models.Session, error) {
	ctx := context.Background()
	session, err := s.store.Session(ctx, userID, state)
	if err != nil {
		return models.Session{}, fmt.Errorf("service: %w", err)
	}
	return session, nil
}
