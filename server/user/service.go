package user

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/0xForked/goca/server/hof"
	"golang.org/x/oauth2"
)

type IUserService interface {
	Profile(ctx context.Context, username string, withPassword bool) (*User, error)
	Availability(ctx context.Context, uid int) (*Availability, error)
	EventType(ctx context.Context, uid int, uname string) ([]*EventType, error)
	SaveGoogleToken(ctx context.Context, username string, googleToken *oauth2.Token) error
	SaveMicrosoftToken(ctx context.Context, username string, microsoftToken *oauth2.Token) error
	Login(ctx context.Context, form *LoginForm) (map[string]interface{}, error)
	Booking(ctx context.Context, uid int) (*Booking, error)
	NewBooking(ctx context.Context, userID int, title string, form *BookingForm, event interface{}) (int, error)
}

type service struct {
	repository ISQLRepository
}

func (s service) Profile(
	ctx context.Context,
	username string,
	withPassword bool,
) (*User, error) {
	user, err := s.repository.FindUserProfile(ctx, username)
	if err != nil {
		return nil, err
	}
	if !withPassword {
		user.Password = ""
	}
	return user, nil
}

func (s service) Availability(
	ctx context.Context,
	uid int,
) (*Availability, error) {
	return s.repository.FindUserAvailability(ctx, uid)
}

func (s service) EventType(
	ctx context.Context,
	uid int,
	uname string,
) ([]*EventType, error) {
	user, err := s.repository.FindUserProfile(ctx, uname)
	if err != nil {
		return nil, err
	}
	eventTypes, err := s.repository.FindUserEventType(ctx, uid)
	if err != nil {
		return nil, err
	}
	for _, eventType := range eventTypes {
		eventType.IsGoogleAvailable = user.GoogleToken.Valid
		eventType.IsMicrosoftAvailable = user.MicrosoftToken.Valid
	}
	return eventTypes, nil
}

func (s service) SaveGoogleToken(
	ctx context.Context,
	username string,
	googleToken *oauth2.Token,
) error {
	data, err := json.Marshal(googleToken)
	if err != nil {
		return err
	}
	return s.repository.UpdateUser(
		ctx, username, "google_token", data)
}

func (s service) SaveMicrosoftToken(
	ctx context.Context,
	username string,
	microsoftToken *oauth2.Token,
) error {
	data, err := json.Marshal(microsoftToken)
	if err != nil {
		return err
	}
	return s.repository.UpdateUser(
		ctx, username, "microsoft_token", data)
}

func (s service) Login(
	ctx context.Context,
	form *LoginForm,
) (map[string]interface{}, error) {
	user, err := s.Profile(ctx, form.Username, true)
	if err != nil {
		return nil, err
	}
	if err := s.validatePassword(user.Password, form.Password); err != nil {
		return nil, err
	}
	at, ate, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"type":       "Bearer",
		"token":      at,
		"expires_in": ate,
	}, nil
}

func (s service) validatePassword(hash, userPwd string) error {
	h := hof.PasswordHash{Stored: hash, Supplied: userPwd}
	isValid, err := h.ComparePassword(hof.Parallelization)
	if err != nil {
		return err
	}
	if !isValid {
		return errors.New("password is not valid")
	}
	return nil
}

func (s service) generateToken(
	user *User,
) (at string, ate *time.Time, err error) {
	jwtToken := hof.JSONWebToken{}
	jwtToken.IssuedAt = time.Now()
	tokenExpiredIn := jwtToken.IssuedAt.Add(time.Minute * 30)
	jwtToken.ExpiredAt = tokenExpiredIn
	accessToken, err := jwtToken.Claim(map[string]interface{}{
		"id":       user.ID,
		"username": user.Username,
	})
	if err != nil {
		return "", nil, err
	}
	return accessToken, &tokenExpiredIn, nil
}

func (s service) Booking(ctx context.Context, uid int) (*Booking, error) {
	return s.repository.FindBooking(ctx, uid)
}

func (s service) NewBooking(
	ctx context.Context,
	userID int,
	title string,
	form *BookingForm,
	event interface{},
) (int, error) {
	newEvent, err := json.Marshal(event)
	if err != nil {
		return 0, err
	}
	newBooking := Booking{
		UserID:      userID,
		EventTypeID: form.EventTypeID,
		Title:       title,
		Notes:       form.Notes,
		Name:        form.Name,
		Email:       form.Email,
		Date:        form.Date,
		Time:        form.Time,
		Event:       newEvent,
		Location:    form.MeetingLocation,
	}
	return s.repository.InsertBooking(ctx, &newBooking)
}

func newUserService(
	repository ISQLRepository,
) IUserService {
	return &service{repository: repository}
}
