package data

import (
	"context"
	"database/sql"
	"github.com/Skaifai/gophers-online-store/internal/validator"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"os"
	"testing"
	"time"
)

func TestValidateProduct(t *testing.T) {
	v := validator.New()

	noNameProduct := &Product{
		Name:        "",
		Price:       100.0,
		Description: "SomeDescription",
		Category:    "Category",
		Quantity:    5,
	}

	ValidateProduct(v, noNameProduct)

	expected := v.Valid()

	if expected != false {
		t.Errorf("ValidateProduct(noNameProduct) returned unexpected value: got %v, expected %s", expected, "false")
	}
}

func TestTableDrivenValidateProduct(t *testing.T) {
	noNameProduct := &Product{
		Name:        "",
		Price:       100.0,
		Description: "SomeDescription",
		Category:    "Category",
		Quantity:    5,
	}

	longNameProduct := &Product{
		Name:        "This A Very Long Name That Has More Than Twenty Bytes In It",
		Price:       100.0,
		Description: "SomeDescription",
		Category:    "Category",
		Quantity:    5,
	}

	badPriceProduct := &Product{
		Name:        "GoodName",
		Price:       -100.0,
		Description: "SomeDescription",
		Category:    "Category",
		Quantity:    5,
	}

	emptyDescriptionProduct := &Product{
		Name:        "GoodName",
		Price:       100.0,
		Description: "",
		Category:    "Category",
		Quantity:    5,
	}

	emptyCategoryProduct := &Product{
		Name:        "GoodName",
		Price:       100.0,
		Description: "Descrpition",
		Category:    "",
		Quantity:    5,
	}

	badQuantityProduct := &Product{
		Name:        "GoodName",
		Price:       100.0,
		Description: "",
		Category:    "Category",
		Quantity:    -1,
	}

	perfectProduct := &Product{
		Name:        "GoodName",
		Price:       100.0,
		Description: "SomeDescription",
		Category:    "Category",
		Quantity:    5,
	}

	var tests = []struct {
		name     string
		input    *Product
		expected bool
	}{
		{
			"ValidateProduct(noNameProduct) must return false",
			noNameProduct,
			false,
		},
		{
			"ValidateProduct(longNameProduct) must return false",
			longNameProduct,
			false,
		},
		{
			"ValidateProduct(badPriceProduct) must return false",
			badPriceProduct,
			false,
		},
		{
			"ValidateProduct(emptyDescriptionProduct) must return false",
			emptyDescriptionProduct,
			false,
		},
		{
			"ValidateProduct(emptyCategoryProduct) must return false",
			emptyCategoryProduct,
			false,
		},
		{
			"ValidateProduct(badQuantityProduct) must return false",
			badQuantityProduct,
			false,
		},
		{
			"ValidateProduct(perfectProduct) must return true",
			perfectProduct,
			true,
		},
	}

	for _, tst := range tests {
		v := validator.New()
		t.Run(tst.name, func(t *testing.T) {
			ValidateProduct(v, tst.input)
			result := v.Valid()
			if result != tst.expected {
				t.Errorf("Expected %v got %v", tst.expected, result)
			}
		})
	}
}

func TestDBConnection(t *testing.T) {
	db, err := sql.Open("postgres", getEnvVar("DB_URL"))

	if err != nil {
		t.Errorf("Database could not be connected!")
	}

	db.SetMaxIdleConns(25)
	db.SetMaxOpenConns(25)

	duration, err := time.ParseDuration("15m")
	if err != nil {
		t.Errorf("Time could not be parsed!")
	}
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		t.Errorf("Could not connect to DB!")
	}
	db.Close()
}

func TestSetCurrentStatus(t *testing.T) {
	someProduct := &Product{
		Name:        "Hello",
		Price:       100.0,
		Description: "SomeDescription",
		Category:    "Category",
		Quantity:    5,
	}

	someProduct.SetStatus(someProduct.Quantity)

	expected := someProduct.IsAvailable

	if expected != true {
		t.Errorf("SetStatus(noNameProduct.Quantity) returned unexpected value: got %v, expected %s", expected, "true")
	}
}

func TestSetNewStatus(t *testing.T) {
	someProduct := &Product{
		Name:        "Hello",
		Price:       100.0,
		Description: "SomeDescription",
		Category:    "Category",
		Quantity:    5,
	}

	someProduct.SetStatus(0)

	expected := someProduct.IsAvailable

	if expected != false {
		t.Errorf("SetStatus(noNameProduct.Quantity) returned unexpected value: got %v, expected %s", expected, "false")
	}
}

func getEnvVar(key string) string {
	godotenv.Load("C:\\Users\\skaif\\Documents\\GitHub\\gophers-online-store\\.env")
	return os.Getenv(key)
}
