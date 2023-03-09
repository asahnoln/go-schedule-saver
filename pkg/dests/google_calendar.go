package dests

import (
	"fmt"
	"time"

	"github.com/asahnoln/go-schedule-saver/pkg"
	"google.golang.org/api/calendar/v3"
)

type GoogleCalendar struct {
	calId string
	srv   *calendar.Service
}

func NewGoogleCalendar(calendarId string, srv *calendar.Service) *GoogleCalendar {
	return &GoogleCalendar{
		calId: calendarId,
		srv:   srv,
	}
}

func (c *GoogleCalendar) Save(es []pkg.Event) error {
	for _, e := range es {
		var h, m, d int
		var mon, day string
		fmt.Sscanf(e.Time, "%2d:%2d", &h, &m)
		fmt.Sscanf(e.Day, "%s\n%d %s", &day, &d, &mon)

		startTime := time.Date(time.Now().Year(), translateMonth(mon), d, h, m, 0, 0, time.Local)
		endTime := startTime.Add(30 * time.Minute)
		_, err := c.srv.Events.Insert(c.calId, &calendar.Event{
			Summary: e.Desc,
			Start: &calendar.EventDateTime{
				DateTime: startTime.Format(time.RFC3339),
			},
			End: &calendar.EventDateTime{
				DateTime: endTime.Format(time.RFC3339),
			},
		}).Do()

		if err != nil {
			return fmt.Errorf("dests gcal save: %w", err)
		}
	}
	return nil
}

func translateMonth(mon string) time.Month {
	// TODO: Протестировать на остальные месяцы
	month := time.January
	switch mon {
	case "марта":
		month = time.March
	case "апреля":
		month = time.April
	}

	return month
}
