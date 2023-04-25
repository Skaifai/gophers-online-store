package main

import (
	"context"
	"encoding/json"
	"github.com/Skaifai/gophers-online-store/internal/data"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

const id int64 = 19

func TestRegisterUserHandler(t *testing.T) {
	recorder := httptest.NewRecorder()
	body := `{
		"name": "Dinmukhammed",
		"surname": "Zhasulanov",
		"username": "dinmukhammed",
		"date_of_birth": "2004-07-04T00:00:00Z",
		"phone_number": "+77474938382",
		"address": "Kabanbay Batyr Avenue 60/11",
		"email": "dinmuhammed.017@gmail.com",
		"password": "iamdimok"
	}`
	req, err := http.NewRequest("POST", "/v1/auth/register", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	testApp.registerUserHandler(recorder, req)

	if recorder.Code != http.StatusAccepted {
		t.Errorf("Expected status 202; got %d", recorder.Code)
	}

	var user data.User
	err = json.Unmarshal(recorder.Body.Bytes(), &user)
	if err != nil {
		t.Fatal("Expected data has JSON value, but it doesn't")
	}

	if user.Activated {
		t.Errorf("Expected activation false, but got %v", user.Activated)
	}
}

func TestActivateUserHandler(t *testing.T) {
	uuid, err := testApp.models.ActivationLinks.GetActivationLink(id)
	if err != nil {
		t.Fatal("No need to continue, uuid has not been found")
	}
	params := httprouter.Params{{Key: "uuid", Value: uuid}}
	ctx := context.WithValue(context.Background(), httprouter.ParamsKey, params)

	recorder := httptest.NewRecorder()
	req, err := http.NewRequestWithContext(ctx, "GET", "/v1/auth/activate/:uuid", nil)
	if err != nil {
		t.Fatal(err)
	}

	testApp.activateUserHandler(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status 200; got %d", recorder.Code)
	}

	var envelope map[string]data.User
	err = json.Unmarshal(recorder.Body.Bytes(), &envelope)
	if err != nil {
		t.Fatalf("Unexpected user found, %v", envelope)
	}

	user := envelope["user"]
	if !user.Activated {
		t.Errorf("Expected activation true, but got %v", user.Activated)
	}
}

func TestAuthenticateUserHandler(t *testing.T) {
	recorder := httptest.NewRecorder()
	body := `{"email": "dinmuhammed.017@gmail.com", "password": "iamdimok"}`
	req, err := http.NewRequest("POST", "/v1/auth/authenticate", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	testApp.authenticateUserHandler(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status 200; got %d", recorder.Code)
	}

	res := recorder.Body.String()
	if !strings.Contains(res, `"accessToken"`) || !strings.Contains(res, `"refreshToken"`) {
		t.Errorf("Unexpected body: %s", res)
	}
}

func TestShowUserHandler(t *testing.T) {
	params := httprouter.Params{{Key: "id", Value: strconv.FormatInt(id, 10)}}
	ctx := context.WithValue(context.Background(), httprouter.ParamsKey, params)

	req, err := http.NewRequestWithContext(ctx, "GET", "/v1/users/:id", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	testApp.showUserHandler(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status 200; got %d", recorder.Code)
	}

	var envelope map[string]data.User
	err = json.Unmarshal(recorder.Body.Bytes(), &envelope)
	if err != nil {
		t.Fatalf("Unexpected user found, %v", envelope)
	}

	user := envelope["user"]
	if user.Username != "dinmukhammed" || user.Email != "dinmuhammed.017@gmail.com" {
		t.Errorf("Unexpected user found, %v", user)
	}
}

func TestUpdateUserHandler(t *testing.T) {
	params := httprouter.Params{{Key: "id", Value: strconv.FormatInt(id, 10)}}
	ctx := context.WithValue(context.Background(), httprouter.ParamsKey, params)
	body := `{"username": "dimok"}`

	req, err := http.NewRequestWithContext(ctx, "PATCH", "/v1/users/:id", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	testApp.updateUserHandler(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status 200; got %d", recorder.Code)
	}

	var envelope map[string]data.User
	err = json.Unmarshal(recorder.Body.Bytes(), &envelope)
	if err != nil {
		t.Fatalf("Unexpected user found, %v", envelope)
	}

	user := envelope["user"]
	if user.Username != "dimok" || user.Email != "dinmuhammed.017@gmail.com" {
		t.Errorf("Unexpected user found, %v", user)
	}
}

func TestDeleteUserHandler(t *testing.T) {
	params := httprouter.Params{{Key: "id", Value: strconv.FormatInt(id, 10)}}
	ctx := context.WithValue(context.Background(), httprouter.ParamsKey, params)

	req, err := http.NewRequestWithContext(ctx, "DELETE", "/v1/users/:id", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	testApp.deleteUserHandler(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status 200; got %d", recorder.Code)
	}

	expectedBody := `{"message":"user successfully deleted"}` + "\n"
	if recorder.Body.String() != expectedBody {
		t.Errorf("Unexpected body: %s", recorder.Body.String())
	}
}
