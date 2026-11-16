package service

import (
	"errors"
	"testing"

	"hw5/internal/model"
)

var errRepository = errors.New("repository error")

type repositoryStub struct {
	createdOrder   *model.Order
	updatedOrder   *model.Order
	deletedID      int
	requestedID    int
	customer       string
	createErr      error
	getByIDErr     error
	updateErr      error
	deleteErr      error
	getItemsErr    error
	getCustomerErr error
}

func (r *repositoryStub) Create(order *model.Order) (*model.Order, error) {
	r.createdOrder = order
	return order, r.createErr
}

func (r *repositoryStub) GetByID(id int) (*model.Order, error) {
	r.requestedID = id
	return &model.Order{ID: id}, r.getByIDErr
}

func (r *repositoryStub) Update(order *model.Order) (*model.Order, error) {
	r.updatedOrder = order
	return order, r.updateErr
}

func (r *repositoryStub) Delete(id int) error {
	r.deletedID = id
	return r.deleteErr
}

func (r *repositoryStub) GetItems(id int) ([]model.Item, error) {
	r.requestedID = id
	return []model.Item{{Name: "coffee", Price: 250}}, r.getItemsErr
}

func (r *repositoryStub) GetByCustomer(customer string) ([]*model.Order, error) {
	r.customer = customer
	return []*model.Order{{ID: 1, Customer: customer}}, r.getCustomerErr
}

func TestCreateAssignsIDAndStatus(t *testing.T) {
	repository := &repositoryStub{}
	service := New(repository)
	input := &model.Order{
		Customer: "Alice",
		Status:   "paid",
		Items:    []model.Item{{Name: "coffee", Price: 250}},
	}

	created, err := service.Create(input)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if created.ID != 1 || created.Status != "new" {
		t.Fatalf("Create() = %#v, want ID 1 and status new", created)
	}
	if input.ID != 0 || input.Status != "paid" {
		t.Fatal("Create() changed the input order")
	}
}

func TestCreateUsesSequentialIDs(t *testing.T) {
	repository := &repositoryStub{}
	service := New(repository)

	first, err := service.Create(validOrder())
	if err != nil {
		t.Fatalf("first Create() error = %v", err)
	}
	second, err := service.Create(validOrder())
	if err != nil {
		t.Fatalf("second Create() error = %v", err)
	}
	if first.ID != 1 || second.ID != 2 {
		t.Fatalf("created IDs = %d, %d; want 1, 2", first.ID, second.ID)
	}
}

func TestCreateValidatesOrder(t *testing.T) {
	tests := []struct {
		name  string
		order *model.Order
	}{
		{name: "nil order"},
		{name: "empty customer", order: &model.Order{Items: []model.Item{{Name: "coffee", Price: 250}}}},
		{name: "empty items", order: &model.Order{Customer: "Alice"}},
		{name: "empty item name", order: &model.Order{Customer: "Alice", Items: []model.Item{{Price: 250}}}},
		{name: "invalid item price", order: &model.Order{Customer: "Alice", Items: []model.Item{{Name: "coffee"}}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := &repositoryStub{}
			_, err := New(repository).Create(tt.order)
			if !errors.Is(err, ErrInvalidOrder) {
				t.Fatalf("Create() error = %v, want %v", err, ErrInvalidOrder)
			}
			if repository.createdOrder != nil {
				t.Fatal("Create() called repository for invalid order")
			}
		})
	}
}

func TestUpdateValidatesAndCallsRepository(t *testing.T) {
	repository := &repositoryStub{}
	service := New(repository)
	order := validOrder()
	order.ID = 3
	order.Status = "paid"

	updated, err := service.Update(order)
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if updated.ID != 3 || repository.updatedOrder == nil {
		t.Fatalf("Update() = %#v, repository got %#v", updated, repository.updatedOrder)
	}

	order.Status = "unknown"
	if _, err := service.Update(order); !errors.Is(err, ErrInvalidOrder) {
		t.Fatalf("Update() invalid status error = %v, want %v", err, ErrInvalidOrder)
	}
}

func TestReadAndDeleteOperations(t *testing.T) {
	repository := &repositoryStub{}
	service := New(repository)

	if _, err := service.GetByID(4); err != nil || repository.requestedID != 4 {
		t.Fatalf("GetByID() error = %v, requested ID = %d", err, repository.requestedID)
	}
	if _, err := service.GetItems(5); err != nil || repository.requestedID != 5 {
		t.Fatalf("GetItems() error = %v, requested ID = %d", err, repository.requestedID)
	}
	if err := service.Delete(6); err != nil || repository.deletedID != 6 {
		t.Fatalf("Delete() error = %v, deleted ID = %d", err, repository.deletedID)
	}
	if _, err := service.GetByCustomer("Alice"); err != nil || repository.customer != "Alice" {
		t.Fatalf("GetByCustomer() error = %v, customer = %q", err, repository.customer)
	}
}

func TestServiceWrapsRepositoryErrors(t *testing.T) {
	repository := &repositoryStub{
		getByIDErr:     errRepository,
		updateErr:      errRepository,
		deleteErr:      errRepository,
		getItemsErr:    errRepository,
		getCustomerErr: errRepository,
	}
	service := New(repository)
	order := validOrder()
	order.ID = 1
	order.Status = "new"

	operations := []struct {
		name string
		run  func() error
	}{
		{name: "get by id", run: func() error { _, err := service.GetByID(1); return err }},
		{name: "update", run: func() error { _, err := service.Update(order); return err }},
		{name: "delete", run: func() error { return service.Delete(1) }},
		{name: "get items", run: func() error { _, err := service.GetItems(1); return err }},
		{name: "get by customer", run: func() error { _, err := service.GetByCustomer("Alice"); return err }},
	}

	for _, operation := range operations {
		t.Run(operation.name, func(t *testing.T) {
			if err := operation.run(); !errors.Is(err, errRepository) {
				t.Fatalf("error = %v, want wrapped repository error", err)
			}
		})
	}
}

func validOrder() *model.Order {
	return &model.Order{
		Customer: "Alice",
		Items:    []model.Item{{Name: "coffee", Price: 250}},
	}
}
