package model

import (
	"database/sql"

	"github.com/golodash/galidator"
)

type User struct {
	ID       int            `json:"id"`
	Username string         `json:"username"`
	Password string         `json:"password,omitempty"`
	Token    sql.NullString `json:"-"`
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
