package data

import (
	"context"
	"database/sql"
	"time"
)

type ShoppingSession struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	Total        float64   `json:"total"`
	CreationDate time.Time `json:"creation_date"`
}

type ShoppingSessionModel struct {
	DB *sql.DB
}

func (s ShoppingSessionModel) Insert(userID int64) error {
	query := `
	INSERT INTO shopping_session (user_id, total) 
	VALUES ($1, $2)
	RETURNING version`

	args := []any{userID, 0}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var version int64
	err := s.DB.QueryRowContext(ctx, query, args...).Scan(&version)

	if err != nil {
		return err
	}
	return nil
}

func (s ShoppingSessionModel) Get(userID int64) (*ShoppingSession, error) {
	query := `SELECT id, user_id, total, creation_date FROM carts WHERE user_id = $1`

	var session ShoppingSession

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := s.DB.QueryRowContext(ctx, query, userID).Scan(
		&session.ID,
		&session.UserID,
		&session.Total,
		&session.CreationDate,
	)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (s ShoppingSessionModel) GetCartItems(id int64) ([]*CartItem, error) {
	query := `
	SELECT id, session_id, product_id, quantity, creation_date
	FROM cart_items
	WHERE session_id = $1`

	var items []*CartItem

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := s.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var item CartItem
		err := rows.Scan(
			&item.ID,
			&item.SessionID,
			&item.ProductID,
			&item.Quantity,
			&item.CreationDate,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, &item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (s ShoppingSessionModel) GetCartTotal(id int64) (float64, error) {
	query := `
	SELECT SUM(price * quantity) 
	FROM cart_items 
	INNER JOIN products ON cart_items.product_id = products.id 
	WHERE session_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var total float64
	err := s.DB.QueryRowContext(ctx, query, id).Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}