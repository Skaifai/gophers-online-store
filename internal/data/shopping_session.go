package data

import (
	"database/sql"
	"time"
)

type ShoppingSession struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	total        float64   `json:"total"`
	CreationDate time.Time `json:"creation_date"`
}

type ShoppingSessionModel struct {
	DB *sql.DB
}
