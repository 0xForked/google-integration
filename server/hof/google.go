package hof

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/people/v1"
)

func GetProfileData(
	svr *people.Service,
) (string, error) {
	profile, err := svr.People.Get("people/me").
		PersonFields("names,emailAddresses").Do()
	if err != nil {
		return "", fmt.Errorf("unable to retrieve user profile: %v", err)
	}
	if len(profile.Names) < 0 {
		return "", nil
	}
	return profile.Names[0].DisplayName, nil
}

func GetPeopleService(
	ctx context.Context,
	tok *oauth2.Token,
	config *oauth2.Config,
) *people.Service {
	client := config.Client(ctx, tok)
	srv, err := people.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}
	return srv
}

func GetCalendarData(
	svr *calendar.Service,
) ([]*calendar.Event, error) {
	t := time.Now().Format(time.RFC3339)
	events, err := svr.Events.List("primary").
		ShowDeleted(false).
		SingleEvents(true).
		TimeMin(t).
		MaxResults(10).
		OrderBy("startTime").
		Do()
	if err != nil {
		return nil, fmt.Errorf(
			"unable to retrieve next ten of the user's events: %v",
			err)
	}
	return events.Items, nil
}

func GetCalendarService(
	ctx context.Context,
	tok *oauth2.Token,
	config *oauth2.Config,
) *calendar.Service {
	client := config.Client(ctx, tok)
	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}
	return srv
}

func GetOAuthConfig() *oauth2.Config {
	b, err := os.ReadFile("google.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	config, err := google.ConfigFromJSON(b,
		calendar.CalendarReadonlyScope,
		people.UserinfoProfileScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	return config
}

func GetOAuthTokenFromWeb(config *oauth2.Config) (string, *oauth2.Token) {
	authURL := config.AuthCodeURL(
		"state-token", oauth2.AccessTypeOffline)
	if authURL != "" {
		return authURL, nil
	}
	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}
	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return "", tok
}
