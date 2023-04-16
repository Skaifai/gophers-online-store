package data

import "testing"

func TestSortDirectionDesc(t *testing.T) {
	newFilter := &Filters{
		Page:     1,
		PageSize: 5,
		Sort:     "-",
		SortSafelist: []string{"id", "name", "category", "price", "is_available", "creation_date",
			"-id", "-name", "-category", "-price", "-is_available", "-creation_date"},
	}

	expected := newFilter.sortDirection()

	if expected != "DESC" {
		t.Errorf("sortDirection() returned unexpected value: got %v, expected %s", expected, "DESC")
	}
}

func TestSortDirectionAsc(t *testing.T) {
	newFilter := &Filters{
		Page:     1,
		PageSize: 5,
		Sort:     "",
		SortSafelist: []string{"id", "name", "category", "price", "is_available", "creation_date",
			"-id", "-name", "-category", "-price", "-is_available", "-creation_date"},
	}

	expected := newFilter.sortDirection()

	if expected != "ASC" {
		t.Errorf("sortDirection() returned unexpected value: got %v, expected %s", expected, "ASC")
	}
}
