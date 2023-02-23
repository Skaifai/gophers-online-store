package data

import (
	"database/sql"
	"time"
)

type CartItem struct {
	ID           int64     `json:"id"`
	SessionID    int64     `json:"user_id"`
	ProductID    int64     `json:"product_id"`
	Quantity     int64     `json:"quantity"`
	CreationDate time.Time `json:"creation_date"`
}

type CartItemModel struct {
	DB *sql.DB
}
