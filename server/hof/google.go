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
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
	"google.golang.org/api/people/v1"
)

func GetUserDisplayName(
	ctx context.Context,
	token *oauth2.Token,
	config *oauth2.Config,
) (string, error) {
	client := config.Client(ctx, token)
	srv, err := people.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return "", fmt.Errorf("Unable to retrieve Calendar client: %v", err)
	}
	profile, err := srv.People.Get("people/me").
		PersonFields("names,emailAddresses").Do()
	if err != nil {
		return "", fmt.Errorf("unable to retrieve user profile: %v", err)
	}
	if len(profile.Names) < 0 {
		return "", nil
	}
	return profile.Names[0].DisplayName, nil
}

func GetUserEmail(
	ctx context.Context,
	token *oauth2.Token,
	config *oauth2.Config,
) (string, error) {
	client := config.Client(ctx, token)
	svc, err := gmail.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return "", fmt.Errorf("unable to create Gmail service: %v", err)
	}
	profile, err := svc.Users.GetProfile("me").Do()
	if err != nil {
		return "", fmt.Errorf("unable to retrieve user's profile: %v", err)
	}
	return profile.EmailAddress, nil
}

func SetNewMeeting(
	svr *calendar.Service,
	summary, description,
	timezone, oEmail, cEmail string,
	date int64, timeInt, duration int,
) (*calendar.Event, error) {
	randStr, err := generateRandomString(12)
	if err != nil {
		log.Fatalf("Unable to generate random string: %v", err)
	}
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		log.Fatalf("Unable to load timezone: %v", err)
	}
	dateObj := time.Unix(date, 0).In(loc)
	dateObj = time.Date(dateObj.Year(), dateObj.Month(), dateObj.Day(), 0, 0, 0, 0, loc)
	dateObj = dateObj.Add(time.Duration(timeInt/100) * time.Hour)
	dateObj = dateObj.Add(time.Duration(timeInt%100) * time.Minute)
	endObj := dateObj.Add(time.Duration(duration) * time.Minute)
	event := &calendar.Event{
		Summary:     summary,
		Description: description,
		Start: &calendar.EventDateTime{
			DateTime: dateObj.Format(time.RFC3339),
			TimeZone: loc.String(),
		},
		End: &calendar.EventDateTime{
			DateTime: endObj.Format(time.RFC3339),
			TimeZone: loc.String(),
		},
		Attendees: []*calendar.EventAttendee{
			{Email: oEmail, Organizer: true, ResponseStatus: "accepted"},
			{Email: cEmail, ResponseStatus: "accepted"},
		},
		ConferenceData: &calendar.ConferenceData{
			CreateRequest: &calendar.CreateConferenceRequest{
				RequestId: randStr,
			},
		},
	}
	return svr.Events.Insert("primary", event).
		ConferenceDataVersion(1).Do()
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
		calendar.CalendarEventsScope,
		people.UserinfoProfileScope,
		gmail.GmailReadonlyScope)
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
