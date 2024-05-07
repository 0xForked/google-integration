package user

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/0xForked/goca/server/hof"
	"github.com/0xForked/goca/server/model"
	"golang.org/x/oauth2"
)

type IUserService interface {
	Profile(ctx context.Context, username string, withPassword bool) (*model.User, error)
	SaveGoogleToken(ctx context.Context, username string, googleToken *oauth2.Token) error
	Login(ctx context.Context, form *model.LoginForm) (map[string]interface{}, error)
}

type service struct {
	repository ISQLRepository
}

func (s service) Profile(
	ctx context.Context,
	username string,
	withPassword bool,
) (*model.User, error) {
	user, err := s.repository.Find(ctx, username)
	if err != nil {
		return nil, err
	}
	if !withPassword {
		user.Password = ""
	}
	return user, nil
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
	return s.repository.Update(ctx, username, data)
}

func (s service) Login(
	ctx context.Context,
	form *model.LoginForm,
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
	user *model.User,
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

func newUserService(
	repository ISQLRepository,
) IUserService {
	return &service{repository: repository}
}
