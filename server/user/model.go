package user

import (
	"database/sql"

	"github.com/golodash/galidator"
)

type User struct {
	ID             int            `json:"id"`
	Username       string         `json:"username"`
	Password       string         `json:"password,omitempty"`
	GoogleToken    sql.NullString `json:"-"`
	MicrosoftToken sql.NullString `json:"-"`
	Availability   *Availability  `json:"availability,omitempty"`
	EventTypes     []*EventType   `json:"event_types,omitempty"`
}

type Availability struct {
	ID       int                `json:"id"`
	UserID   int                `json:"-"`
	Label    string             `json:"label"`
	Timezone string             `json:"timezone"`
	Days     []*AvailabilityDay `json:"days"`
}

type AvailabilityDay struct {
	ID             int `json:"id"`
	UserID         int `json:"-"`
	AvailabilityID int `json:"-"`
	Enable         int `json:"enable"`
	Day            int `json:"day"`        // store it 0-6
	StartTime      int `json:"start_time"` // use TimeToInt func
	EndTime        int `json:"end_time"`   // use TimeToInt func
}

type EventType struct {
	ID                   int           `json:"id"`
	UserID               int           `json:"-"`
	AvailabilityID       int           `json:"-"`
	Enable               int           `json:"enable"`
	Title                string        `json:"title"`
	Description          string        `json:"description"`
	Duration             int           `json:"duration"` // duration in minute
	Availability         *Availability `json:"availability,omitempty"`
	IsGoogleAvailable    bool          `json:"is_google_available"`
	IsMicrosoftAvailable bool          `json:"is_microsoft_available"`
}

type Booking struct {
	ID          string      `json:"id"`
	UserID      int         `json:"-"`
	EventTypeID int         `json:"-"`
	Title       string      `json:"title"`
	Notes       string      `json:"notes"`
	Name        string      `json:"name"`
	Email       string      `json:"email"`
	Date        int64       `json:"date"`
	Time        int         `json:"time"`
	Location    string      `json:"location"`
	Event       []byte      `json:"-"`
	EventDetail interface{} `json:"event_detail"`
}

type LoginForm struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

func (f *LoginForm) Validate() interface{} {
	g := galidator.New()
	return g.ComplexValidator(galidator.Rules{
		"Username": g.R("username").Required(),
		"Password": g.R("password").Required(),
	}).Validate(f)
}

type BookingForm struct {
	Username        string `json:"username" form:"username"`
	EventTypeID     int    `json:"event_type_id" form:"event_type_id"`
	Date            int64  `json:"date" form:"date"`
	Time            int    `json:"time" form:"time"`
	Name            string `json:"name" form:"name"`
	Email           string `json:"email" form:"email"`
	Notes           string `json:"notes" form:"notes"`
	MeetingLocation string `json:"meeting_location" form:"meeting_location"`
}

func (f *BookingForm) Validate() interface{} {
	g := galidator.New()
	return g.ComplexValidator(galidator.Rules{
		"Username":    g.R("username").Required(),
		"EventTypeID": g.R("event_type_id").Required(),
		"Date":        g.R("date").Required(),
		"Time":        g.R("time").Required(),
		"Name":        g.R("name").Required(),
		"Email":       g.R("email").Required(),
	}).Validate(f)
}
