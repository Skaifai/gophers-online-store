package data

import (
	"database/sql"
)

type Token struct {
	ID           int64  `json:"-"`
	UserID       int64  `json:"-"`
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token,omitempty"`
}

type TokenModel struct {
	DB *sql.DB
}
