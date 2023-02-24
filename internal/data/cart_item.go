package data

import (
	"context"
	"database/sql"
	"errors"
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

func (i CartItemModel) Get(id int64) (*CartItem, error) {
	query := `
	SELECT session_id, product_id, quantity, creation_date
	FROM cart_items
	WHERE id = $1`

	var item CartItem
	item.ID = id

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := i.DB.QueryRowContext(ctx, query, id).Scan(
		&item.SessionID,
		&item.ProductID,
		&item.Quantity,
		&item.CreationDate,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &item, nil
}

func (i CartItemModel) Insert(item *CartItem) error {
	query := `
	INSERT INTO cart_items (session_id, product_id, quantity)
	VALUES ($1, $2, $3)
	RETURNING id`

	args := []any{item.SessionID, item.ProductID, item.Quantity}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := i.DB.QueryRowContext(ctx, query, args...).Scan(&item.ID)
	if err != nil {
		return err
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

	query := `DELETE FROM cart_items WHERE id = $1`
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
