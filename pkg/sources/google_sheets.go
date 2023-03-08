package sources

import (
	"github.com/asahnoln/go-schedule-saver/pkg"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/sheets/v4"
)

const (
	headerRow       = 0
	timeColumn      = 0
	valsColumnStart = 1
	valsRowStart    = 1
)

type Doer interface {
	Do(opts ...googleapi.CallOption) (*sheets.ValueRange, error)
}

type GoogleSheets struct {
	d Doer
}

func NewGoogleSheets(d Doer) *GoogleSheets {
	return &GoogleSheets{d}
}

func (s *GoogleSheets) Parse() ([]pkg.Event, error) {
	resp, _ := s.d.Do()
	events := []pkg.Event{}

	dates := resp.Values[headerRow]
	for _, v := range resp.Values[valsRowStart:] {
		t := v[timeColumn]

		for c := valsColumnStart; c < len(dates); c++ {
			desc := v[c].(string)
			if desc == "" {
				continue
			}

			events = append(events, pkg.Event{
				Day:  dates[c].(string),
				Time: t.(string),
				Desc: desc,
			})
		}
	}

	return events, nil
}
