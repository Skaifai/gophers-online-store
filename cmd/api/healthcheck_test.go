package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthcheckHandler(t *testing.T) {
	recorder := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/v1/healthcheck", nil)
	if err != nil {
		t.Fatal(err)
	}

	testApp.healthcheckHandler(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status 200; got %d", recorder.Code)
	}

	expectedBody := `{"status":"available","system_info":{"environment":"test","version":"1.0"}}` + "\n"
	if recorder.Body.String() != expectedBody {
		t.Errorf("Unexpected body: %s", recorder.Body.String())
	}
}
