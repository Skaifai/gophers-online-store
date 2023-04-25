package main

import (
	"context"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

//func TestListCommentsHandler(t *testing.T) {
//	cfg := SetupConfig()
//
//	db, err := openDB(cfg)
//	if err != nil {
//		t.Error(err)
//	}
//	defer db.Close()
//
//	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)
//
//	app := SetupApplication(cfg, logger, db)
//
//	server := httptest.NewServer(http.HandlerFunc(app.listCommentsHandler))
//
//	resp, err := http.Get(server.URL)
//	if err != nil {
//		t.Error(err)
//	}
//
//	if resp.StatusCode != http.StatusOK {
//		t.Errorf("Expected 200, got %d", resp.StatusCode)
//	}
//
//	server.Close()
//}
func TestListCommentsHandler(t *testing.T) {
	params := httprouter.Params{{Key: "id", Value: strconv.FormatInt(id, 10)}}
	ctx := context.WithValue(context.Background(), httprouter.ParamsKey, params)

	req, err := http.NewRequestWithContext(ctx, "GET", "/v1/products/:id/comments", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	testingApplication.listCommentsHandler(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status 200; got %d", recorder.Code)
	}
}
