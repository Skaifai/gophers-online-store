package data

import (
	"database/sql"
	"time"
)

type Product struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Price        float64   `json:"price"`
	Description  string    `json:"description"`
	Category     string    `json:"category"`
	Quantity     int       `json:"-"`
	IsAvailable  bool      `json:"is_avalable"`
	CreationDate time.Time `json:"creation_date"`
	Version      int       `json:"-"`
}

type ProductModel struct {
	DB *sql.DB
}
