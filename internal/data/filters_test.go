package data

import "testing"

var filterByNameDesc = Filters{
	Page:         5,
	PageSize:     10,
	Sort:         "-name",
	SortSafelist: []string{"id", "-id", "name", "-name"},
}

var filterByIdAsc = Filters{
	Page:         2,
	PageSize:     5,
	Sort:         "id",
	SortSafelist: []string{"id", "-id", "name", "-name"},
}

func TestSortColumn(t *testing.T) {
	column := filterByNameDesc.sortColumn()
	expected := "name"
	if column != expected {
		t.Errorf("sortColumn() returned %s, expected %s", column, expected)
	}

	column = filterByIdAsc.sortColumn()
	expected = "id"
	if column != expected {
		t.Errorf("sortColumn() returned %s, expected %s", column, expected)
	}
}

func TestSortDirection(t *testing.T) {
	expected := "DESC"
	direction := filterByNameDesc.sortDirection()
	if direction != expected {
		t.Errorf("sortDirection() returned %s, expected %s", direction, expected)
	}

	expected = "ASC"
	direction = filterByIdAsc.sortDirection()
	if direction != expected {
		t.Errorf("sortDirection() returned %s, expected %s", direction, expected)
	}
}

func TestLimit(t *testing.T) {
	expected := 10
	limit := filterByNameDesc.limit()
	if limit != expected {
		t.Errorf("limit() returned %d, expected %d", limit, expected)
	}

	expected = 5
	limit = filterByIdAsc.limit()
	if limit != expected {
		t.Errorf("limit() returned %d, expected %d", limit, expected)
	}
}

func TestOffSet(t *testing.T) {
	expected := 40
	offset := filterByNameDesc.offset()
	if offset != expected {
		t.Errorf("offset() returned %d, expected %d", offset, expected)
	}

	expected = 5
	offset = filterByIdAsc.offset()
	if offset != expected {
		t.Errorf("offset() returned %d, expected %d", offset, expected)
	}
}
