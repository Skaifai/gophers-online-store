package data

import (
	"database/sql"
	"time"
)

type Comment struct {
	ID           int64     `json:"id"`
	ProductID    int64     `json:"product_id"`
	CommentOwner int64     `json:"owner_id"`
	Text         string    `json:"text"`
	CreationDate time.Time `json:"creation_date"`
}

type CommentModel struct {
	DB *sql.DB
}
