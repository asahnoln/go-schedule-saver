package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/asahnoln/go-schedule-saver/pkg"
	"github.com/asahnoln/go-schedule-saver/pkg/dests"
	"github.com/asahnoln/go-schedule-saver/pkg/sources"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
func destLogger(e []pkg.Event) error {
	fmt.Printf("Events:\n%v\n", e)
	return nil
}

func readMailsFile() map[string]string {
	f, err := os.Open(os.Getenv("MAILS_PATH"))
	if err != nil {
		log.Fatalf("Read mails error: %s", err)
	}

	var mails map[string]string
	err = json.NewDecoder(f).Decode(&mails)
	if err != nil {
		log.Fatalf("Decode mails error: %s", err)
	}

	return mails
}

func createCalendar() *calendar.Service {
	ctx := context.Background()
	b, err := os.ReadFile(os.Getenv("CAL_CREDS"))
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, calendar.CalendarScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	return srv
}

func main() {
	srv, err := sheets.NewService(context.Background(), option.WithCredentialsFile(os.Getenv("SHEETS_CREDS")))
	if err != nil {
		log.Fatalf("Sheets creation error: %s", err)
	}
	// cal, err := calendar.NewService(context.Background(), option.WithCredentialsFile(os.Getenv("CAL_CREDS")))
	// if err != nil {
	// 	log.Fatalf("Calendar creation error: %s", err)
	// }

	s := sources.NewGoogleSheets(srv.Spreadsheets.Values.Get(os.Getenv("SHEETS_ID"), os.Getenv("SHEETS_RANGE")))
	d := dests.NewGoogleCalendar(os.Getenv("CAL_ID"), createCalendar(), readMailsFile())
	err = pkg.Save(s, d)
	if err != nil {
		log.Fatalf("Save error: %s", err)
	}
}
