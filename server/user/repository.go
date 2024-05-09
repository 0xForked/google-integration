package user

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
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
	Update(
		ctx context.Context,
		username string,
		token []byte,
	) error
}

type sqlRepository struct {
	db *sql.DB
}

//goland:noinspection ALL
func (s sqlRepository) FindUserProfile(
	ctx context.Context,
	username string,
) (*User, error) {
	q := "SELECT id, username, password, google_token FROM users WHERE username = ?"
	row := s.db.QueryRowContext(ctx, q, username)
	var user User
	err := row.Scan(&user.ID, &user.Username,
		&user.Password, &user.Token)
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
func (s sqlRepository) Update(
	ctx context.Context,
	username string,
	token []byte,
) error {
	q := "UPDATE users SET google_token = ? WHERE username = ?"
	_, err := s.db.ExecContext(ctx, q, token, username)
	if err != nil {
		return err
	}
	return nil
}

func newSQLRepository(db *sql.DB) ISQLRepository {
	return sqlRepository{db: db}
}
