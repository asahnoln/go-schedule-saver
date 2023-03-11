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
	calEvents, _ := c.srv.Events.List(c.calId).Do()
	esMap := make(map[string]bool)
	if calEvents != nil {
		for _, e := range calEvents.Items {
			esMap[e.Summary+e.Start.DateTime] = true
		}
	}

	for _, e := range es {
		var h, m, d int
		var mon, day string
		fmt.Sscanf(e.Time, "%2d:%2d", &h, &m)
		fmt.Sscanf(e.Day, "%s\n%d %s", &day, &d, &mon)

		startTime := time.Date(time.Now().Year(), translateMonth(mon), d, h, m, 0, 0, time.Local)

		if esMap[e.Desc+startTime.Format(time.RFC3339)] {
			continue
		}

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
