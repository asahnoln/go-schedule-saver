package dests_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/asahnoln/go-schedule-saver/pkg"
	"github.com/asahnoln/go-schedule-saver/pkg/dests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func TestGoogleCalendarSavesEvents(t *testing.T) {
	got := []calendar.Event{}
	fs := fakeServer(t, &got, "calId")
	s, err := calendar.NewService(context.Background(),
		option.WithoutAuthentication(),
		option.WithEndpoint(fs.URL))
	require.NoError(t, err)

	c := dests.NewGoogleCalendar("calId", s, make(map[string]string))
	require.Implements(t, (*pkg.Destination)(nil), c)

	want := []pkg.Event{
		{Day: "Понедельник,\n2 марта", Time: "14:30", Desc: "Terry"},
		{Day: "Понедельник,\n3 марта", Time: "16:00", Desc: "Mike"},
		{Day: "Среда,\n4 апреля", Time: "15:30", Desc: "Lola"},
	}
	err = c.Save(want)
	require.NoError(t, err)

	assert.Len(t, got, len(want))
	assert.Equal(t, want[0].Desc, got[0].Summary)
	assert.Equal(t,
		time.Date(time.Now().Year(), time.March, 2, 14, 30, 0, 0, time.Local).Format(time.RFC3339),
		got[0].Start.DateTime)
	assert.Equal(t,
		time.Date(time.Now().Year(), time.March, 2, 15, 00, 0, 0, time.Local).Format(time.RFC3339),
		got[0].End.DateTime)

	assert.Equal(t,
		time.Date(time.Now().Year(), time.April, 4, 15, 30, 0, 0, time.Local).Format(time.RFC3339),
		got[2].Start.DateTime)
}

func TestErrorFromServer(t *testing.T) {
	fs := errorFakeServer()
	s, _ := calendar.NewService(context.Background(),
		option.WithoutAuthentication(),
		option.WithEndpoint(fs.URL))

	c := dests.NewGoogleCalendar("testCalId", s, make(map[string]string))
	err := c.Save([]pkg.Event{
		{},
	})
	require.Error(t, err)
}

func TestCalendarDoesNotSaveDuplicates(t *testing.T) {
	got := []calendar.Event{
		{
			// FIX: Actually it shouldn't be saved, it should be moved to a different slot
			Summary: "Terry", // Different Terry event
			Start: &calendar.EventDateTime{
				DateTime: time.Date(time.Now().Year(), time.May, 1, 12, 15, 0, 0, time.Local).Format(time.RFC3339),
			},
		},
		{
			Summary: "Dora", // Already exists, shouldn't be duplicated
			Start: &calendar.EventDateTime{
				DateTime: time.Date(time.Now().Year(), time.April, 4, 15, 30, 0, 0, time.Local).Format(time.RFC3339),
			},
		},
	}
	fs := fakeServer(t, &got, "dupCalId")
	s, _ := calendar.NewService(context.Background(),
		option.WithoutAuthentication(),
		option.WithEndpoint(fs.URL))

	c := dests.NewGoogleCalendar("dupCalId", s, make(map[string]string))

	want := []pkg.Event{
		{Day: "Понедельник,\n2 марта", Time: "14:30", Desc: "Terry"},
		{Day: "Понедельник,\n3 марта", Time: "16:00", Desc: "Mike"},
		{Day: "Среда,\n4 апреля", Time: "15:30", Desc: "Dora"},
	}
	err := c.Save(want)
	require.NoError(t, err)

	assert.Len(t, got, 4)
}

// TODO: Test for error from list

func TestAddsGuestsAccordingToGivenNames(t *testing.T) {
	got := []calendar.Event{}
	fs := fakeServer(t, &got, "guestsCalId")
	s, _ := calendar.NewService(context.Background(),
		option.WithoutAuthentication(),
		option.WithEndpoint(fs.URL))

	mails := map[string]string{
		"Alex": "alex@mail.com",
		"Bob":  "bob@mail.com",
	}
	c := dests.NewGoogleCalendar("guestsCalId", s, mails)

	want := []pkg.Event{
		{Day: "Понедельник,\n2 мая", Time: "11:30", Desc: "Alex"},
		{Day: "Понедельник,\n3 мая", Time: "12:00", Desc: "Bob"},
	}
	err := c.Save(want)
	require.NoError(t, err)

	assert.Equal(t, "alex@mail.com", got[0].Attendees[0].Email)
}

// TODO: Separate fake server concerns (listing, empty events, etc)
func fakeServer(t *testing.T, got *[]calendar.Event, calId string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.Split(r.URL.Path, "/")
		require.Equal(t, calId, path[2], "google calendar id is wrong")

		// List events
		if r.Method == http.MethodGet {
			g := *got
			items := []*calendar.Event{}
			for i := range g {
				items = append(items, &g[i])
			}
			es := &calendar.Events{
				Items: items,
			}

			resp, _ := es.MarshalJSON()
			w.Write(resp)
			return
		}

		// Post events
		var e calendar.Event
		err := json.NewDecoder(r.Body).Decode(&e)
		defer r.Body.Close()

		if err != nil {
			t.Fatal(err)
		}

		if len(e.Attendees) > 0 && e.Attendees[0].Email == "" {
			t.Fatalf("Attendee email cannot be empty")
		}

		*got = append(*got, e)

		_, _ = w.Write([]byte("{}"))
	}))
}

func errorFakeServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "some error", http.StatusInternalServerError)
	}))
}
