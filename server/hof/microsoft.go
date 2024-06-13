package hof

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
)

type MSEvent struct {
	Subject               string          `json:"subject"`
	Body                  MSBody          `json:"body"`
	Start                 MSEventStartEnd `json:"start"`
	End                   MSEventStartEnd `json:"end"`
	Attendees             []MSAttendee    `json:"attendees"`
	IsOnlineMeeting       bool            `json:"isOnlineMeeting"`
	OnlineMeetingProvider string          `json:"onlineMeetingProvider"`
}

type MSBody struct {
	ContentType string `json:"contentType"`
	Content     string `json:"content"`
}

type MSEventStartEnd struct {
	DateTime string `json:"dateTime"`
	TimeZone string `json:"timeZone"`
}

type MSAttendee struct {
	EmailAddress MSEmailAddress `json:"emailAddress"`
	Type         string         `json:"type"`
}

type MSEmailAddress struct {
	Address string `json:"address"`
	Name    string `json:"name"`
}

func ComposeMSMeetingData(
	timezone, summary string,
	date int64, timeInt, duration int,
	menteeName, menteeMail string,
) MSEvent {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		log.Fatalf("Unable to load timezone: %v", err)
	}
	dateObj := time.Unix(date, 0).In(loc)
	dateObj = time.Date(dateObj.Year(), dateObj.Month(), dateObj.Day(), 0, 0, 0, 0, loc)
	dateObj = dateObj.Add(time.Duration(timeInt/100) * time.Hour)
	dateObj = dateObj.Add(time.Duration(timeInt%100) * time.Minute)
	endObj := dateObj.Add(time.Duration(duration) * time.Minute)
	return MSEvent{
		Subject: summary,
		Start: MSEventStartEnd{
			DateTime: dateObj.Format(time.RFC3339),
			TimeZone: loc.String(),
		},
		End: MSEventStartEnd{
			DateTime: endObj.Format(time.RFC3339),
			TimeZone: loc.String(),
		},
		Attendees: []MSAttendee{
			{
				EmailAddress: MSEmailAddress{
					Address: menteeMail,
					Name:    menteeName,
				},
				Type: "required",
			},
		},
		IsOnlineMeeting:       true,
		OnlineMeetingProvider: "teamsForBusiness",
	}
}

func SetMicrosoftNewCalendarEvent(
	event MSEvent,
	accessToken string,
) (map[string]interface{}, error) {
	// Convert event to JSON
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return nil, fmt.Errorf("error marshalling event to JSON: %s", err.Error())
	}
	// Make request to Microsoft Graph API to create event
	url := "https://graph.microsoft.com/v1.0/me/calendar/events"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(eventJSON))
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %s", err.Error())
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making HTTP request:%s", err.Error())
	}
	defer func() { _ = resp.Body.Close() }()
	// Parse the response to get the meeting link
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %s", err.Error())
	}
	var eventData map[string]interface{}
	if err := json.Unmarshal(responseBody, &eventData); err != nil {
		return nil, fmt.Errorf("error unmarshalling response body: %s", err.Error())
	}
	return eventData, nil
}

func SetMicrosoftNewMeeting(
	startDate, endDate, subject, accessToken string,
) (map[string]interface{}, error) {
	// Convert event to JSON
	eventJSON, err := json.Marshal(map[string]interface{}{
		"startDateTime": startDate,
		"endDateTime":   endDate,
		"subject":       subject,
	})
	if err != nil {
		return nil, fmt.Errorf("error marshalling event to JSON: %s", err.Error())
	}
	// Make request to Microsoft Graph API to create event
	url := "https://graph.microsoft.com/v1.0/me/onlineMeetings"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(eventJSON))
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %s", err.Error())
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making HTTP request:%s", err.Error())
	}
	defer func() { _ = resp.Body.Close() }()
	// Parse the response to get the meeting link
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %s", err.Error())
	}
	var meetingData map[string]interface{}
	if err := json.Unmarshal(responseBody, &meetingData); err != nil {
		return nil, fmt.Errorf("error unmarshalling response body: %s", err.Error())
	}
	return meetingData, nil
}

func GetMicrosoftUserProfile(accessToken string) (map[string]string, error) {
	url := "https://graph.microsoft.com/v1.0/me"
	req, err := http.NewRequest("GET", url, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("error creating http request: %s", err.Error())
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making http request: %s", err.Error())
	}
	defer func() { _ = resp.Body.Close() }()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %s", err.Error())
	}
	var userData map[string]interface{}
	if err := json.Unmarshal(responseBody, &userData); err != nil {
		return nil, fmt.Errorf("error unmarshalling response body: %s", err.Error())
	}
	if userData == nil {
		return nil, nil
	}
	return map[string]string{
		"name": func() string {
			if name, ok := userData["displayName"].(string); ok {
				return name
			}
			return ""
		}(),
		"email": func() string {
			if mail, ok := userData["mail"].(string); ok {
				return mail
			}
			return ""
		}(),
	}, err
}

func GetMicrosoftOAuthConfig() *oauth2.Config {
	b, err := os.ReadFile("microsoft.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	type microsoftCredentials struct {
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		//TenantID     string `json:"tenant_id"`
		RedirectURI string `json:"redirect_uri"`
	}
	var j struct {
		Web *microsoftCredentials `json:"web"`
	}
	if err := json.Unmarshal(b, &j); err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	return &oauth2.Config{
		ClientID:     j.Web.ClientID,
		ClientSecret: j.Web.ClientSecret,
		RedirectURL:  j.Web.RedirectURI,
		Endpoint:     microsoft.AzureADEndpoint(""),
		Scopes: []string{
			"User.Read", "email", "openid", "profile", "offline_access",
			"Calendars.Read", "Calendars.ReadWrite",
			"OnlineMeetings.ReadWrite",
		},
	}
}

func GetMicrosoftOAuthTokenFromWeb(
	config *oauth2.Config,
) (string, *oauth2.Token) {
	authURL := config.AuthCodeURL(
		"state", oauth2.AccessTypeOffline)
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

//wt-ug: https://github.com/calcom/cal.com/blob/33d7da88bfda9375c17c4302a0c27b9f64a15d5d/packages/app-store/office365calendar/lib/CalendarService.ts
//ms-doc: https://learn.microsoft.com/en-us/graph/api/calendar-post-events?view=graph-rest-1.0&tabs=http
//https://learn.microsoft.com/en-us/graph/api/application-post-onlinemeetings?view=graph-rest-1.0&tabs=http
