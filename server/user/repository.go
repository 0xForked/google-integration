package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/0xForked/goca/server/model"
)

type ISQLRepository interface {
	Find(
		ctx context.Context,
		username string,
	) (*model.User, error)
	Update(
		ctx context.Context,
		username string,
		token []byte,
	) error
}

type sqlRepository struct {
	db *sql.DB
}

func (s sqlRepository) Find(
	ctx context.Context,
	username string,
) (*model.User, error) {
	//goland:noinspection ALL
	q := "SELECT id, username, password, google_token FROM users WHERE username = ?"
	row := s.db.QueryRowContext(ctx, q, username)
	var user model.User
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

func (s sqlRepository) Update(
	ctx context.Context,
	username string,
	token []byte,
) error {
	//goland:noinspection ALL
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
