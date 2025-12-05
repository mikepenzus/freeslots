package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type CalendarExporterStatus struct {
	CredentialsFileName     string
	TokenFileName           string
	WebserverAddressAndPort string
}

func CreateCalendarService(calendarExporterStatus CalendarExporterStatus) (*calendar.Service, error) {
	ctx := context.Background()

	// Read credentials from JSON file
	credentialsFile, err := os.ReadFile(calendarExporterStatus.CredentialsFileName)
	if err != nil {
		return nil, err
	}

	// Configure OAuth2 with required scopes
	oauthConfiguration, err := google.ConfigFromJSON(credentialsFile, calendar.CalendarReadonlyScope)
	if err != nil {
		return nil, err
	}
	oauthConfiguration.RedirectURL = "http://" + calendarExporterStatus.WebserverAddressAndPort

	client := getOAuthClient(oauthConfiguration, calendarExporterStatus)
	if err != nil {
		return nil, err
	}

	// Create Calendar service
	calendarService, err := calendar.NewService(ctx, option.WithHTTPClient(client))

	return calendarService, err
}

// getOAuthClient retrieves a token, saves it, then returns the configured client
func getOAuthClient(config *oauth2.Config, calendarExporterStatus CalendarExporterStatus) *http.Client {
	tok, err := tokenFromFile(calendarExporterStatus.TokenFileName)
	if err != nil {
		tok = getTokenFromWeb(config, calendarExporterStatus.WebserverAddressAndPort)
		saveToken(calendarExporterStatus.TokenFileName, tok)
	}
	return config.Client(context.Background(), tok)
}

// getTokenFromWeb requests a token using a local web server
func getTokenFromWeb(config *oauth2.Config, webserverAddressAndPort string) *oauth2.Token {
	codeCh := make(chan string)
	errCh := make(chan error)

	// Start local server to receive OAuth callback
	server := &http.Server{Addr: webserverAddressAndPort}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			errCh <- fmt.Errorf("no code in response")
			return
		}

		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, "<html><body><h1>Authentication successful!</h1><p>You can close this window and return to the terminal.</p></body></html>")
		codeCh <- code
	})

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	// Generate auth URL with localhost redirect
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Opening browser for authentication...\n")
	fmt.Printf("If the browser doesn't open, go to:\n%v\n\n", authURL)

	// Try to open browser (this may not work on all systems)
	fmt.Println("Waiting for authentication...")

	// Wait for code or error
	var authCode string
	select {
	case authCode = <-codeCh:
		// Success
	case err := <-errCh:
		log.Fatalf("Error during authentication: %v", err)
	case <-time.After(5 * time.Minute):
		log.Fatal("Authentication timeout")
	}

	// Shutdown server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(ctx)

	tok, err := config.Exchange(context.Background(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token: %v", err)
	}
	return tok
}

// tokenFromFile retrieves a token from a local file
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

// saveToken saves a token to a file path
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.Create(path)
	if err != nil {
		log.Fatalf("Unable to cache token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
