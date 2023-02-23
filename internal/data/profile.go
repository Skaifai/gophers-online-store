package data

import (
	"database/sql"
	"time"
)

type Profile struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	Surname    string    `json:"surname"`
	DOB        time.Time `json:"date_of_birth"`
	Address    string    `json:"address"`
	AboutMe    string    `json:"about_me"`
	PictureURL string    `json:"picture_url"`
}

type ProfileModel struct {
	DB *sql.DB
}
