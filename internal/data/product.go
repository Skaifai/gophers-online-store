package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Skaifai/gophers-online-store/internal/validator"
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

func (p *Product) SetStatus(productQuantity int) {
	if productQuantity > 0 {
		p.IsAvailable = true
	} else {
		p.IsAvailable = false
	}
}

func ValidateProduct(v *validator.Validator, product *Product) {
	v.Check(product.Name != "", "name", "must be provided")
	v.Check(len(product.Name) <= 20, "name", "must not be more than 20 bytes long")
	v.Check(product.Price >= 0, "price", "can not be negative")
	v.Check(product.Description != "", "description", "must be provided")
	v.Check(product.Category != "", "category", "must be provided")
	v.Check(product.Quantity >= 0, "quantity", "can not be negative")
}

func (p ProductModel) Insert(product *Product) error {
	query := `
	INSERT INTO products (name, price, description, category, quantity, is_available)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, creation_date, version`

	args := []any{
		product.Name,
		product.Price,
		product.Description,
		product.Category,
		product.Quantity,
		product.IsAvailable,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := p.DB.QueryRowContext(ctx, query, args...).Scan(&product.ID,
		&product.CreationDate, &product.Version)

	if err != nil {
		return err
	}

	return nil
}

func (p ProductModel) GetAll(name string, category string, filters Filters) ([]*Product, Metadata, error) {
	// Update the SQL query to include the window function which counts the total
	// (filtered) records.
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), id, name, price, description, category, is_available, creation_date, version
		FROM products
		WHERE (to_tsvector('simple', name) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (to_tsvector('simple', category) @@ plainto_tsquery('simple', $2) OR $2 = '')
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	args := []any{name, category, filters.limit(), filters.offset()}

	rows, err := p.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err // Update this to return an empty Metadata struct.
	}

	defer rows.Close()

	totalRecords := 0

	var products []*Product

	for rows.Next() {
		var product Product
		err := rows.Scan(
			&totalRecords, // Scan the count from the window function into totalRecords.
			&product.ID,
			&product.Name,
			&product.Price,
			&product.Description,
			&product.Category,
			&product.IsAvailable,
			&product.CreationDate,
			&product.Version,
		)
		if err != nil {
			return nil, Metadata{}, err // Update this to return an empty Metadata struct.
		}
		products = append(products, &product)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err // Update this to return an empty Metadata struct.
	}
	// Generate a Metadata struct, passing in the total record count and pagination
	// parameters from the client.
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	// Include the metadata struct when returning.
	return products, metadata, nil
}

func (p ProductModel) Get(id int64) (*Product, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	// Define the SQL query for retrieving the movie data.
	query := `
		SELECT id, name, price, description, category, is_available, creation_date, version
		FROM products
		WHERE id = $1`
	// Declare a Movie struct to hold the data returned by the query.
	var product Product
	err := p.DB.QueryRow(query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Price,
		&product.Description,
		&product.Category,
		&product.IsAvailable,
		&product.CreationDate,
		&product.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &product, nil
}

func (p ProductModel) Update(product *Product) error {
	query := `UPDATE products
	SET name = $1, price = $2, description = $3, category = $4, quantity = $5, version = version + 1
	WHERE id = $6
	RETURNING version`

	args := []any{
		product.Name,
		product.Price,
		product.Description,
		product.Category,
		product.Quantity,
		product.ID,
	}

	return p.DB.QueryRow(query, args...).Scan(&product.Version)
}

func (p ProductModel) Delete(id int64) error {
	query := `
		DELETE FROM products
		WHERE id = $1`
	result, err := p.DB.Exec(query, id)
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
