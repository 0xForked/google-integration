package user

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"time"
)

type ISQLRepository interface {
	FindUserProfile(
		ctx context.Context,
		username string,
	) (*User, error)
	FindUserAvailability(
		ctx context.Context,
		uid int,
	) (*Availability, error)
	FindUserEventType(
		ctx context.Context,
		uid int,
	) ([]*EventType, error)
	UpdateUser(
		ctx context.Context,
		username string,
		tokenType string,
		tokenValue []byte,
	) error
	FindBooking(
		ctx context.Context,
		bookingID int,
	) (*Booking, error)
	InsertBooking(
		ctx context.Context,
		booking *Booking,
	) (int, error)
}

type sqlRepository struct {
	db *sql.DB
}

//goland:noinspection ALL
func (s sqlRepository) FindUserProfile(
	ctx context.Context,
	username string,
) (*User, error) {
	q := "SELECT id, username, password, google_token, microsoft_token "
	q += "FROM users WHERE username = ?"
	row := s.db.QueryRowContext(ctx, q, username)
	var user User
	err := row.Scan(&user.ID, &user.Username,
		&user.Password, &user.GoogleToken, &user.MicrosoftToken)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf(
				"account with username %s not found",
				username)
		}
		return nil, err
	}
	return &user, nil
}

//goland:noinspection ALL
func (s sqlRepository) FindUserAvailability(
	ctx context.Context,
	uid int,
) (*Availability, error) {
	q := `SELECT
	    a.id,
	    a.user_id,
	    a.label,
	    a.timezone,
	    (SELECT json_group_array(
			json_object(
				'id', ad.id,
				'user_id', ad.user_id,
				'availability_id', ad.availability_id,
				'Enable', ad.Enable,
				'day', ad.day,
				'start_time', ad.start_time,
				'end_time', ad.end_time
			))
	        FROM availability_days AS ad
	        WHERE ad.availability_id = a.id
	    ) AS availability_days
	FROM availabilities AS a WHERE a.user_id = 1`
	row := s.db.QueryRowContext(ctx, q, uid)
	var av Availability
	var availabilityDaysJSON []byte
	if err := row.Scan(&av.ID, &av.UserID,
		&av.Label, &av.Timezone, &availabilityDaysJSON); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf(
				"user availabililty with id %d not found",
				uid)
		}
		return nil, err
	}
	var availabilityDays []*AvailabilityDay
	if err := json.Unmarshal(availabilityDaysJSON, &availabilityDays); err != nil {
		return nil, fmt.Errorf("failed to unmarshal availability_days: %v", err)
	}
	av.Days = availabilityDays
	return &av, nil
}

//goland:noinspection ALL
func (s sqlRepository) FindUserEventType(
	ctx context.Context,
	uid int,
) ([]*EventType, error) {
	q := `SELECT
	    et.id,
	    et.user_id,
	    et.availability_id,
	    et.enable,
	    et.title,
	    et.description,
	    et.duration,
	    a.id as av_id,
	    a.label as av_label,
	    a.timezone as av_timezone,
        (
            SELECT json_group_array(
                json_object(
                    'id', ad.id,
                    'user_id', ad.user_id,
                    'availability_id', ad.availability_id,
                    'Enable', ad.Enable,
                    'day', ad.day,
                    'start_time', ad.start_time,
                    'end_time', ad.end_time
                )
            )
            FROM availability_days AS ad
            WHERE ad.availability_id = a.id
        ) AS availability_days
	FROM event_types AS et
	LEFT JOIN availabilities AS a ON et.availability_id = a.id
	WHERE et.user_id = ?`
	rows, err := s.db.QueryContext(ctx, q, uid)
	if err != nil {
		return nil, err
	}
	var eventTypes []*EventType
	for rows.Next() {
		var et EventType
		var av Availability
		var availabilityDaysJSON []byte
		if err := rows.Scan(
			&et.ID, &et.UserID, &et.AvailabilityID,
			&et.Enable, &et.Title, &et.Description,
			&et.Duration, &av.ID, &av.Label, &av.Timezone,
			&availabilityDaysJSON,
		); err != nil {
			return nil, err
		}
		var availabilityDays []*AvailabilityDay
		if err := json.Unmarshal(availabilityDaysJSON, &availabilityDays); err != nil {
			return nil, fmt.Errorf("failed to unmarshal availability_days: %v", err)
		}
		sort.Slice(availabilityDays, func(i, j int) bool {
			return availabilityDays[i].Enable > availabilityDays[j].Enable
		})
		et.Availability = &Availability{
			ID:       av.ID,
			Label:    av.Label,
			Timezone: av.Timezone,
			Days:     availabilityDays,
		}
		eventTypes = append(eventTypes, &et)
	}
	return eventTypes, nil
}

//goland:noinspection ALL
func (s sqlRepository) UpdateUser(
	ctx context.Context,
	username string,
	tokenType string,
	tokenValue []byte,
) error {
	// google_token || microsoft_token
	q := "UPDATE users SET "
	q += fmt.Sprintf("%s = ? WHERE username = ?", tokenType)
	_, err := s.db.ExecContext(ctx, q, tokenValue, username)
	if err != nil {
		return err
	}
	return nil
}

//goland:noinspection ALL
func (s sqlRepository) FindBooking(
	ctx context.Context,
	bookingID int,
) (*Booking, error) {
	q := "SELECT id, title, notes, name, email, date, time, event FROM bookings WHERE id = ?"
	row := s.db.QueryRowContext(ctx, q, bookingID)
	var booking Booking
	var bookingJSON []byte
	err := row.Scan(&booking.ID, &booking.Title, &booking.Notes,
		&booking.Name, &booking.Email, &booking.Date, &booking.Time, &bookingJSON)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf(
				"booking with id %d not found",
				bookingID)
		}
		return nil, err
	}
	if err := json.Unmarshal(bookingJSON, &booking.EventDetail); err != nil {
		return nil, fmt.Errorf("failed to unmarshal availability_days: %v", err)
	}
	return &booking, nil
}

//goland:noinspection ALL
func (s sqlRepository) InsertBooking(
	ctx context.Context,
	booking *Booking,
) (int, error) {
	q := "INSERT INTO bookings (user_id, event_type_id, title, notes, name, email, date, time, event, location, created_at) "
	q += "values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id"
	row := s.db.QueryRowContext(ctx, q, booking.UserID, booking.EventTypeID, booking.Title,
		booking.Notes, booking.Name, booking.Email, booking.Date, booking.Time,
		booking.Event, booking.Location, time.Now().Unix())
	var id int
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func newSQLRepository(db *sql.DB) ISQLRepository {
	return sqlRepository{db: db}
}
