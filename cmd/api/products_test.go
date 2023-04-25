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

var id int64 = 10

func TestListProductsHandler(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(testingApplication.listProductsHandler))

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}
	//fmt.Println(resp.Body)
	server.Close()
}

func TestAddProductHandler(t *testing.T) {
	var input = `
		{
			"name": "Football",
			"price": 4000,
			"description": "description4",
			"category": "Sports",
			"quantity": 0
		}`

	server := httptest.NewServer(http.HandlerFunc(testingApplication.addProductHandler))
	defer server.Close()
	//req, err := http.NewRequest("POST", server.URL, strings.NewReader(string(body)))
	//if err != nil {
	//	t.Fatal(err)
	//}

	resp, err := http.Post(server.URL, "application/json", strings.NewReader(input))
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusAccepted {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}
}

func TestShowProductHandler(t *testing.T) {
	params := httprouter.Params{{Key: "id", Value: strconv.FormatInt(id, 10)}}
	ctx := context.WithValue(context.Background(), httprouter.ParamsKey, params)

	req, err := http.NewRequestWithContext(ctx, "GET", "/v1/products/:id", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	testingApplication.showProductHandler(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status 200; got %d", recorder.Code)
	}

	var envelope map[string]data.Product
	err = json.Unmarshal(recorder.Body.Bytes(), &envelope)
	if err != nil {
		t.Fatalf("Couldn't unmarshal")
	}

	product := envelope["product"]
	if product.Name != "Football" || product.Category != "Sports" {
		t.Errorf("Unexpected product found, %v", product)
	}
}

func TestUpdateProductHandler(t *testing.T) {
	var input = `
		{
			"name": "Basketball",
			"price": 3000,
			"description": "description4",
			"category": "Sports",
			"quantity": 0
		}`

	params := httprouter.Params{{Key: "id", Value: strconv.FormatInt(id, 10)}}
	ctx := context.WithValue(context.Background(), httprouter.ParamsKey, params)

	req, err := http.NewRequestWithContext(ctx, "PATCH", "/v1/products/:id", strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	testingApplication.updateProductHandler(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status 200; got %d", recorder.Code)
	}

	var envelope map[string]data.Product
	err = json.Unmarshal(recorder.Body.Bytes(), &envelope)
	if err != nil {
		t.Fatalf("Couldn't unmarshal")
	}

	product := envelope["product"]
	if product.Name != "Basketball" || product.Category != "Sports" {
		t.Errorf("Unexpected product found, %v", product)
	}
}

func TestDeleteProductHandler(t *testing.T) {
	params := httprouter.Params{{Key: "id", Value: strconv.FormatInt(id, 10)}}
	ctx := context.WithValue(context.Background(), httprouter.ParamsKey, params)

	req, err := http.NewRequestWithContext(ctx, "DELETE", "/v1/products/:id", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	testingApplication.deleteProductHandler(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status 200; got %d", recorder.Code)
	}

	if recorder.Body.String() != `{"message":"product successfully deleted"}`+"\n" {
		t.Errorf("Unexpected response; got %s", recorder.Body.String())
	}

	var envelope envelope
	err = json.Unmarshal(recorder.Body.Bytes(), &envelope)
	if err != nil {
		t.Fatalf("Couldn't unmarshal")
	}
	responseMessage := envelope["message"]
	if responseMessage != "product successfully deleted" {
		t.Errorf("Unexpected response; got %s", recorder.Body.String())
	}
}
