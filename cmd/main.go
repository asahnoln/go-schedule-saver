package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/asahnoln/go-schedule-saver/pkg"
	"github.com/asahnoln/go-schedule-saver/pkg/dests"
	"github.com/asahnoln/go-schedule-saver/pkg/sources"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func destLogger(e []pkg.Event) error {
	fmt.Printf("Events:\n%v\n", e)
	return nil
}

func main() {
	srv, err := sheets.NewService(context.Background(), option.WithCredentialsFile(os.Getenv("SHEETS_CREDS")))
	if err != nil {
		log.Fatalf("Sheets creation error: %s", err)
	}
	cal, err := calendar.NewService(context.Background(), option.WithCredentialsFile(os.Getenv("CAL_CREDS")))
	if err != nil {
		log.Fatalf("Calendar creation error: %s", err)
	}

	s := sources.NewGoogleSheets(srv.Spreadsheets.Values.Get(os.Getenv("SHEETS_ID"), os.Getenv("SHEETS_RANGE")))
	d := dests.NewGoogleCalendar(os.Getenv("CAL_ID"), cal)
	err = pkg.Save(s, d)
	if err != nil {
		log.Fatalf("Save error: %s", err)
	}
}
