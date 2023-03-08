package sources

import "github.com/asahnoln/go-schedule-saver/pkg"

type GoogleSheets struct {
}

func NewGoogleSheets() *GoogleSheets {
	return nil
}

func (s *GoogleSheets) Parse() ([]pkg.Event, error) {
	return nil, nil
}
