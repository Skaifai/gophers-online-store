package main

import (
	"context"
	"encoding/json"
	"github.com/Skaifai/gophers-online-store/internal/data"
	_ "github.com/Skaifai/gophers-online-store/internal/data"
	_ "github.com/Skaifai/gophers-online-store/internal/jsonlog"
	_ "github.com/Skaifai/gophers-online-store/internal/validator"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"net/http/httptest"
	_ "os"
	"strconv"
	"testing"
)

func TestShowUserHandler(t *testing.T) {
	var id int64 = 1
	params := httprouter.Params{{Key: "id", Value: strconv.FormatInt(id, 10)}}
	ctx := context.WithValue(context.Background(), httprouter.ParamsKey, params)

	req, err := http.NewRequestWithContext(ctx, "GET", "/v1/users/:id", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	testingApplication.showUserHandler(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status 200; got %d", recorder.Code)
	}

	var envelope map[string]data.User
	err = json.Unmarshal(recorder.Body.Bytes(), &envelope)
	if err != nil {
		t.Fatalf("Unexpected user found, %v", envelope)
	}

	user := envelope["user"]
	if user.Username != "Skaifai" || user.Email != "zakhep82@gmail.com" {
		t.Errorf("Unexpected user found, %v", user)
	}
}
