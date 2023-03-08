package sources_test

import (
	"testing"

	"github.com/asahnoln/go-schedule-saver/pkg"
	"github.com/asahnoln/go-schedule-saver/pkg/sources"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/sheets/v4"
)

func TestSpreadsheetParse(t *testing.T) {
	want := []pkg.Event{
		{Day: "Feb 1, Mon", Time: "13:00", Desc: "Steve"},
		{Day: "Feb 1, Mon", Time: "14:00", Desc: "Pete"},
		{Day: "Feb 3, Wed", Time: "13:30", Desc: "Sam"},
	}
	m := &sheetsMock{want}
	s := sources.NewGoogleSheets(m)
	require.Implements(t, (*pkg.Source)(nil), s)

	got, err := s.Parse()
	require.NoError(t, err)
	assert.ElementsMatch(t, want, got)
}

type sheetsMock struct {
	data []pkg.Event
}

func (d *sheetsMock) Do(opts ...googleapi.CallOption) (*sheets.ValueRange, error) {
	return &sheets.ValueRange{
		Values: [][]interface{}{
			{"Time / Days", "Feb 1, Mon", "Feb 3, Wed"},
			{"13:00", "Steve", ""},
			{"13:30", "", "Sam"},
			{"14:00", "Pete", ""},
		},
	}, nil
}
