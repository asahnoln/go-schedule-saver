package sources

import "github.com/asahnoln/go-schedule-saver/pkg"

type GoogleSpreadsheet struct {
}

func NewGoogleSpreadsheet() *GoogleSpreadsheet {
	return nil
}

func (s *GoogleSpreadsheet) Parse() ([]pkg.Event, error) {
	return nil, nil
}
