package main

import (
	"github.com/Skaifai/gophers-online-store/internal/jsonlog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestHealthCheckHandler(t *testing.T) {
	cfg := SetupConfig()

	db, err := openDB(cfg)
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)
	
	app := SetupApplication(cfg, logger, db)

	server := httptest.NewServer(http.HandlerFunc(app.healthcheckHandler))

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}

	server.Close()
}
