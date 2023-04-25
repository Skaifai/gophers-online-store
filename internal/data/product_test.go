package data

import (
	"errors"
	"fmt"
	"github.com/Skaifai/gophers-online-store/internal/validator"
	_ "github.com/lib/pq"
	_ "net/http"
	"testing"
)

var productModel = testProductModel()

func TestValidateProductByPrice(t *testing.T) {
	product := &Product{
		Name:        "Apple",
		Price:       -100,
		Description: "Apple from Almaty city",
		Category:    "Fruit",
		Quantity:    5,
	}
	v := validator.New()
	if ValidateProduct(v, product); v.Valid() {
		t.Error("validation should thrown, because price has negative value")
	}
}

func TestValidateProductByQuantity(t *testing.T) {
	product := &Product{
		Name:        "Apple",
		Price:       850,
		Description: "Apple from Almaty city",
		Category:    "Fruit",
		Quantity:    -1,
	}
	v := validator.New()
	if ValidateProduct(v, product); v.Valid() {
		t.Error("validation should thrown, because quantity has negative value")
	}
}

func TestGetProduct(t *testing.T) {
	product, err := productModel.Get(1)
	if err != nil {
		t.Fatalf("error acquired while accessing table product. %s", err.Error())
	}
	if product.Name != "Product4" && product.Price != 4000 {
		t.Error("returned another product")
	}
}

func TestGetAllProduct(t *testing.T) {
	var input struct {
		Name     string
		Category string
		Filters
	}
	input.Name = ""
	input.Category = ""
	input.Filters.Page = 1
	input.Filters.PageSize = 20
	input.Filters.Sort = "id"
	input.Filters.SortSafelist = []string{"id", "name", "category", "price", "is_available", "creation_date",
		"-id", "-name", "-category", "-price", "-is_available", "-creation_date"}
	products, metadata, err := productModel.GetAll(input.Name, input.Category, input.Filters)
	if err != nil {
		t.Fatalf("error acquired while accessing table product. %s", err.Error())
	}
	if metadata.PageSize != input.PageSize {
		t.Error("method get all products don't working well")
	}
	for _, product := range products {
		fmt.Println(product.Name)
	}
}

func TestAddProduct(t *testing.T) {
	product := &Product{
		Name:        "Apple",
		Price:       850,
		Description: "Apple from Almaty city",
		Category:    "Fruit",
		Quantity:    5,
	}
	product.SetStatus(product.Quantity)
	if !product.IsAvailable {
		t.Error("product should be available")
	}

	v := validator.New()
	if ValidateProduct(v, product); !v.Valid() {
		t.Fatal("dont have correct struct")
	}

	err := productModel.Insert(product)
	if err != nil {
		t.Fatalf("error acquired while inserting product. %s", err.Error())
	}
}

func TestUpdateProduct(t *testing.T) {
	product := &Product{
		ID:          5,
		Name:        "Apple",
		Price:       1200,
		Description: "Apple from Almaty city",
		Category:    "Fruit",
		Quantity:    3,
	}
	err := productModel.Update(product)
	if err != nil {
		t.Fatalf("error acquired while updating product. %s", err.Error())
	}
	result, _ := productModel.Get(product.ID)
	if result.Name != product.Name && result.Price != product.Price {
		t.Error("returned another product")
	}
}

func TestDeleteProduct(t *testing.T) {
	var id int64 = 5
	err := productModel.Delete(id)
	if err != nil {
		t.Fatalf("error acquired while deleting product. %s", err.Error())
	}
	_, err = productModel.Get(id)
	if !errors.Is(err, ErrRecordNotFound) {
		t.Error("product should be deleted, but it is not")
	}
}

func testProductModel() ProductModel {
	db, _ := DBConnection()
	return ProductModel{DB: db}
}
