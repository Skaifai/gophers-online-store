package main

import (
	_ "github.com/Skaifai/gophers-online-store/internal/data"
	"github.com/Skaifai/gophers-online-store/internal/jsonlog"
	_ "github.com/Skaifai/gophers-online-store/internal/validator"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

//func TestRegisterUserHandler(t *testing.T) {
//	newUser := &data.User{
//		Name:        "Name",
//		Surname:     "Surname",
//		Username:    "username",
//		DOB:         time.Now(),
//		PhoneNumber: "+123456789",
//		Address:     "Some address",
//		Email:       "email@gmail.com",
//	}
//	newUser.Password.Set("somePassword")
//
//	newUserJson, err := json.Marshal(newUser)
//	if err != nil {
//		t.Errorf("Password encryption returned an error. \n%v", err)
//		return
//	}
//
//	app := SetupApplication()
//
//	server := httptest.NewServer(http.HandlerFunc(app.registerUserHandler))
//
//	resp, err := http.Post(server.URL, "application/json", strings.NewReader(string(newUserJson)))
//	if err != nil {
//		t.Error(err)
//	}
//
//	if resp.StatusCode != http.StatusOK {
//		t.Errorf("Expected 200, got %d", resp.StatusCode)
//	}
//
//	//req := httptest.NewRequest("POST", "/v1/auth/register", strings.NewReader(string(newUserJson)))
//	//req.Header.Add("Content-Type", "application/json")
//	//
//	//w := http.ResponseWriter()
//	//
//	//res, err :=
//}

//func TestShowUserHandler(t *testing.T) {
//	app := SetupApplication()
//
//	server := httptest.NewServer(http.HandlerFunc(app.showUserHandler))
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

func TestShowUserHandler(t *testing.T) {
	cfg := SetupConfig()

	db, err := openDB(cfg)
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	app := SetupApplication(cfg, logger, db)

	server := httptest.NewServer(http.HandlerFunc(app.showUserHandler))

	resp, err := http.Get(server.URL + "/1")
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}

	server.Close()
}
