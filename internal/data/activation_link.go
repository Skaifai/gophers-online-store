package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type ActivationLink struct {
	Link         string    `json:"link"`
	Activated    bool      `json:"activated"`
	UserID       int64     `json:"user_id"`
	CreationDate time.Time `json:"creation_date"`
}

type ActivationLinkModel struct {
	DB *sql.DB
}

func (a ActivationLinkModel) Insert(user *User, uuidCode string) error {
	query := `
	INSERT INTO activation_links (link, activated, user_id)
	VALUES ($1, $2, $3)
	RETURNING user_id`

	args := []any{uuidCode, user.Activated, user.ID}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := a.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (a ActivationLinkModel) Get(uuid string) (*ActivationLink, error) {
	if len(uuid) < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
	SELECT link, activated, user_id
	FROM activation_links
	WHERE link = $1`

	var activationLink ActivationLink

	err := a.DB.QueryRow(query, uuid).Scan(
		&activationLink.Link,
		&activationLink.Activated,
		&activationLink.UserID,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &activationLink, nil
}
