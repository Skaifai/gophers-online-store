package data

import (
	"context"
	"database/sql"
	"time"
)

type CartItem struct {
	ID           int64     `json:"id"`
	SessionID    int64     `json:"session_id"`
	ProductID    int64     `json:"product_id"`
	Quantity     int64     `json:"quantity"`
	CreationDate time.Time `json:"creation_date"`
}

type CartItemModel struct {
	DB *sql.DB
}

func (i CartItemModel) Insert(sessionId int64, productID int64, quantity int64) error {
	query := `
	INSERT INTO card_items (sessionId, product_id, quantity)
	VALUES ($1, $2, $3)`

	args := []any{sessionId, productID, quantity}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := i.DB.QueryRowContext(ctx, query, args...)
	if err != nil {
		return err.Err()
	}
	return nil
}

func (i CartItemModel) Update(item *CartItem) error {
	query := `UPDATE cart_items
	SET quantity = $1
	WHERE id = $2`

	return i.DB.QueryRow(query, item.Quantity, item.ID).Err()
}

func (i CartItemModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `DELETE FROM card_items WHERE id = $1`
	result, err := i.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}
