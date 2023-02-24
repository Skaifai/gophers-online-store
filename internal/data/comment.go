package data

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Skaifai/gophers-online-store/internal/validator"
	"time"
)

type Comment struct {
	ID           int64     `json:"id"`
	ProductID    int64     `json:"product_id"`
	CommentOwner int64     `json:"owner_id"`
	Text         string    `json:"text"`
	CreationDate time.Time `json:"creation_date"`
	Version      int       `json:"-"`
}

type CommentModel struct {
	DB *sql.DB
}

func ValidateComment(v *validator.Validator, comment *Comment) {
	v.Check(comment.Text != "", "text", "must be provided")
	v.Check(len(comment.Text) <= 100, "text", "must not be more than 100 bytes long")
}

func (c CommentModel) Insert(comment *Comment) error {
	query := `
	INSERT INTO comments (product_id, owner_id, text)
	VALUES ($1, $2, $3)
	RETURNING id, creation_date, version`

	args := []any{
		comment.ProductID,
		comment.CommentOwner,
		comment.Text,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := c.DB.QueryRowContext(ctx, query, args...).Scan(&comment.ID,
		&comment.CreationDate, &comment.Version)

	if err != nil {
		return err
	}

	return nil
}

func (c CommentModel) GetAll(id int64, filters Filters) ([]*Comment, Metadata, error) {
	// Update the SQL query to include the window function which counts the total
	// (filtered) records.
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), id, product_id, owner_id, text, creation_date, version
		FROM comments
		WHERE product_id = $1
		ORDER BY %s %s, id ASC
		LIMIT $2 OFFSET $3`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	args := []any{id, filters.limit(), filters.offset()}

	rows, err := c.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err // Update this to return an empty Metadata struct.
	}

	defer rows.Close()

	totalRecords := 0

	comments := []*Comment{}

	for rows.Next() {
		var comment Comment
		err := rows.Scan(
			&totalRecords, // Scan the count from the window function into totalRecords.
			&comment.ID,
			&comment.ProductID,
			&comment.CommentOwner,
			&comment.Text,
			&comment.CreationDate,
			&comment.Version,
		)
		if err != nil {
			return nil, Metadata{}, err // Update this to return an empty Metadata struct.
		}
		comments = append(comments, &comment)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err // Update this to return an empty Metadata struct.
	}
	// Generate a Metadata struct, passing in the total record count and pagination
	// parameters from the client.
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	// Include the metadata struct when returning.
	return comments, metadata, nil
}

func (c CommentModel) Update(comment *Comment) error {
	query := `UPDATE comments
	SET text = $1, version = version + 1
	WHERE id = $2
	RETURNING version`

	args := []any{
		comment.Text,
		comment.ID,
	}

	return c.DB.QueryRow(query, args...).Scan(&comment.Version)
}

func (c CommentModel) Delete(id int64) error {
	query := `
		DELETE FROM comments
		WHERE id = $1`
	result, err := c.DB.Exec(query, id)
	if err != nil {
		return nil
	}

	// Checking how many rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// Check if the row was in the database before the query
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}
