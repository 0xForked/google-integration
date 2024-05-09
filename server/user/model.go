package user

import (
	"database/sql"

	"github.com/golodash/galidator"
)

type User struct {
	ID           int            `json:"id"`
	Username     string         `json:"username"`
	Password     string         `json:"password,omitempty"`
	Token        sql.NullString `json:"-"`
	Availability *Availability  `json:"availability,omitempty"`
	EventTypes   []*EventType   `json:"event_types,omitempty"`
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
	ID             int           `json:"id"`
	UserID         int           `json:"-"`
	AvailabilityID int           `json:"-"`
	Enable         int           `json:"enable"`
	Title          string        `json:"title"`
	Description    string        `json:"description"`
	Duration       int           `json:"duration"` // duration in minute
	Availability   *Availability `json:"availability,omitempty"`
}

type LoginForm struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

func (r *LoginForm) Validate() interface{} {
	g := galidator.New()
	return g.ComplexValidator(galidator.Rules{
		"Username": g.R("username").Required(),
		"Password": g.R("password").Required(),
	}).Validate(r)
}
