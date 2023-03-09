package dests_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
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
	fs := fakeServer(t, &got)
	s, err := calendar.NewService(context.Background(),
		option.WithoutAuthentication(),
		option.WithEndpoint(fs.URL))
	require.NoError(t, err)

	c := dests.NewGoogleCalendar("testCalId", s)
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

func fakeServer(t *testing.T, got *[]calendar.Event) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var e calendar.Event
		err := json.NewDecoder(r.Body).Decode(&e)
		defer r.Body.Close()

		if err != nil {
			t.Fatal(err)
		}

		*got = append(*got, e)
	}))
}
