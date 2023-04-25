package data

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	_ "net/http"
	"testing"
	"time"
)

func DBConnection() (*sql.DB, error) {
	db, err := sql.Open("postgres", "postgres://postgres:0000@localhost/gophers?sslmode=disable")
	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(25)
	db.SetMaxOpenConns(25)
	duration, _ := time.ParseDuration("15m")
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func TestDBConnection(t *testing.T) {
	_, err := DBConnection()
	if err != nil {
		t.Fatalf("error acquired while connecting database. %s", err.Error())
	}
}
