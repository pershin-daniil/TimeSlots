package calendar

import (
	"context"
	"embed"
	"os"
	"time"

	"github.com/pershin-daniil/TimeSlots/pkg/models"
	"google.golang.org/api/calendar/v3"

	"github.com/sirupsen/logrus"
)

//nolint:typecheck
//go:embed credentials
var credentials embed.FS

type Calendar struct {
	log *logrus.Entry
	srv *calendar.Service
}

func New(ctx context.Context, log *logrus.Logger) *Calendar {
	srv := calendarService(ctx)
	return &Calendar{
		log: log.WithField("module", "calendar"),
		srv: srv,
	}
}

func (c *Calendar) SlotsOneWeek() []models.Event {
	t := time.Now()
	tString := t.Format(time.RFC3339)
	tPlusWeekString := t.Add(7 * 24 * time.Hour).Format(time.RFC3339)
	events, err := c.srv.Events.List(os.Getenv("CALENDAR_ID")).ShowDeleted(false).
		SingleEvents(true).TimeMin(tString).TimeMax(tPlusWeekString).OrderBy("startTime").Do()
	if err != nil {
		c.log.Panicf("Unable to retrieve next ten of the user's events: %v", err)
	}

	if len(events.Items) == 0 {
		return nil
	}

	result := make([]models.Event, 0, len(events.Items))

	for _, item := range events.Items {
		event := models.Event{
			ID:          item.Id,
			Title:       item.Summary,
			Description: item.Description,
			Start:       item.Start.DateTime,
			End:         item.End.DateTime,
			Created:     item.Created,
			Updated:     item.Updated,
			Status:      item.Status,
		}
		result = append(result, event)
	}

	return result
}
