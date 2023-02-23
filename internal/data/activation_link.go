package data

import (
	"database/sql"
	"time"
)

type ActivationLink struct {
	Link         string    `json:"link"`
	Activated    bool      `json:"activated"`
	UserID       int       `json:"user_id"`
	CreationDate time.Time `json:"creation_date"`
}

type ActivationLinkModel struct {
	DB *sql.DB
}
