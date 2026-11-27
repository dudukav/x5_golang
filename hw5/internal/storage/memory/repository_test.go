package memory

import (
	"errors"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hw5/internal/model"
)

func TestCreateAndGetByID(t *testing.T) {
	repository := New()
	order := &model.Order{ID: 1, Customer: "Alice", Status: "new"}

	created, err := repository.Create(order)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	assert.Equal(t, order, created)

	got, err := repository.GetByID(1)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}
	assert.Equal(t, order, got)
	assert.NotSame(t, order, got)
}

func TestCreateAlreadyExists(t *testing.T) {
	repository := New()
	if _, err := repository.Create(&model.Order{ID: 1}); err != nil {
		t.Fatalf("first Create() error = %v", err)
	}

	_, err := repository.Create(&model.Order{ID: 1})
	if !errors.Is(err, ErrAlreadyExists) {
		t.Fatalf("second Create() error = %v, want %v", err, ErrAlreadyExists)
	}
}

func TestGetByIDNotFound(t *testing.T) {
	repository := New()

	_, err := repository.GetByID(1)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("GetByID() error = %v, want %v", err, ErrNotFound)
	}
}

func TestUpdate(t *testing.T) {
	repository := New()
	order := &model.Order{ID: 1, Status: "new"}
	if _, err := repository.Create(order); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	updated := &model.Order{
		ID:     1,
		Status: "paid",
		Items:  []model.Item{{Name: "coffee", Price: 250}},
	}
	got, err := repository.Update(updated)
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	assert.Equal(t, updated, got)
	assert.NotSame(t, updated, got)
}

func TestUpdateNotFound(t *testing.T) {
	repository := New()

	_, err := repository.Update(&model.Order{ID: 1})
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("Update() error = %v, want %v", err, ErrNotFound)
	}
}

func TestDelete(t *testing.T) {
	repository := New()
	if _, err := repository.Create(&model.Order{ID: 1}); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if err := repository.Delete(1); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	if _, err := repository.GetByID(1); !errors.Is(err, ErrNotFound) {
		t.Fatalf("GetByID() after Delete() error = %v, want %v", err, ErrNotFound)
	}
}

func TestDeleteNotFound(t *testing.T) {
	repository := New()

	if err := repository.Delete(1); !errors.Is(err, ErrNotFound) {
		t.Fatalf("Delete() error = %v, want %v", err, ErrNotFound)
	}
}

func TestGetItems(t *testing.T) {
	repository := New()
	order := &model.Order{
		ID: 1,
		Items: []model.Item{
			{Name: "coffee", Price: 250},
			{Name: "cake", Price: 300},
		},
	}
	if _, err := repository.Create(order); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	items, err := repository.GetItems(1)
	if err != nil {
		t.Fatalf("GetItems() error = %v", err)
	}
	if len(items) != 2 || items[0] != order.Items[0] || items[1] != order.Items[1] {
		t.Fatalf("GetItems() = %#v, want %#v", items, order.Items)
	}

	items[0].Name = "changed"
	if order.Items[0].Name == "changed" {
		t.Fatal("GetItems() returned repository's internal slice")
	}
}

func TestGetItemsNotFound(t *testing.T) {
	repository := New()

	_, err := repository.GetItems(1)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("GetItems() error = %v, want %v", err, ErrNotFound)
	}
}

func TestGetByCustomer(t *testing.T) {
	repository := New()
	orders := []*model.Order{
		{ID: 1, Customer: "Alice"},
		{ID: 2, Customer: "Bob"},
		{ID: 3, Customer: "Alice"},
	}
	for _, order := range orders {
		if _, err := repository.Create(order); err != nil {
			t.Fatalf("Create() error = %v", err)
		}
	}

	got, err := repository.GetByCustomer("Alice")
	if err != nil {
		t.Fatalf("GetByCustomer() error = %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("GetByCustomer() returned %d orders, want 2", len(got))
	}
	ids := []int{got[0].ID, got[1].ID}
	sort.Ints(ids)
	if ids[0] != 1 || ids[1] != 3 {
		t.Fatalf("GetByCustomer() IDs = %v, want [1 3]", ids)
	}
}

func TestGetByCustomerNotFound(t *testing.T) {
	repository := New()

	_, err := repository.GetByCustomer("Alice")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("GetByCustomer() error = %v, want %v", err, ErrNotFound)
	}
}

func TestListIsSortedAndIndependent(t *testing.T) {
	repository := New()
	for _, order := range []*model.Order{{ID: 3}, {ID: 1}, {ID: 2}} {
		_, err := repository.Create(order)
		require.NoError(t, err)
	}

	orders, err := repository.List()
	require.NoError(t, err)
	require.Len(t, orders, 3)
	assert.Equal(t, []int{1, 2, 3}, []int{orders[0].ID, orders[1].ID, orders[2].ID})

	orders[0].Status = model.StatusCompleted
	stored, err := repository.GetByID(1)
	require.NoError(t, err)
	assert.NotEqual(t, model.StatusCompleted, stored.Status)
}
