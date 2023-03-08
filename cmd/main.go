package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/asahnoln/go-schedule-saver/pkg"
	"github.com/asahnoln/go-schedule-saver/pkg/sources"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func destLogger(e []pkg.Event) error {
	fmt.Printf("Events:\n%v\n", e)
	return nil
}

func main() {
	ctx := context.Background()
	srv, err := sheets.NewService(ctx, option.WithCredentialsFile(os.Getenv("SHEETS_CREDS")))
	if err != nil {
		log.Fatalf("Service creation error: %s", err)
	}

	s := sources.NewGoogleSheets(srv.Spreadsheets.Values.Get(os.Getenv("SHEETS_ID"), os.Getenv("SHEETS_RANGE")))
	err = pkg.Save(s, pkg.DestinationFunc(destLogger))
	if err != nil {
		log.Fatalf("Save error: %s", err)
	}
}
