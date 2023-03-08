package sources_test

import (
	"testing"

	"github.com/asahnoln/go-schedule-saver/pkg"
	"github.com/asahnoln/go-schedule-saver/pkg/sources"
	"github.com/stretchr/testify/require"
)

func TestSpreadsheetParse(t *testing.T) {
	s := sources.NewGoogleSheets()
	require.Implements(t, (*pkg.Source)(nil), s)
}
