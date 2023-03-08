package pkg_test

import (
	"errors"
	"testing"

	"github.com/asahnoln/go-schedule-saver/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaver(t *testing.T) {
	data := []pkg.Event{
		{"Jan 1, Mon", "12:00", "Ivan"},
		{"Jan 1, Mon", "13:00", "Michael"},
		{"Jan 3, Wed", "12:30", "John"},
	}
	s := &sourceStub{data, nil}
	d := &destStub{}

	err := pkg.Save(s, d)
	require.NoError(t, err)

	assert.Equal(t, data, d.data)
}

func TestSourceError(t *testing.T) {
	s := &sourceStub{err: errors.New("source error")}
	d := &destStub{}

	err := pkg.Save(s, d)
	require.Error(t, err)
}

func TestDestinationError(t *testing.T) {
	s := &sourceStub{}
	d := &destStub{err: errors.New("destination error")}

	err := pkg.Save(s, d)
	require.Error(t, err)
}

type sourceStub struct {
	data []pkg.Event
	err  error
}

func (s *sourceStub) Parse() ([]pkg.Event, error) {
	return s.data, s.err
}

type destStub struct {
	data []pkg.Event
	err  error
}

func (d *destStub) Save(es []pkg.Event) error {
	d.data = es
	return d.err
}
