package transport

import (
	"testing"

	"hw5/internal/model"
)

type serviceStub struct {
	createdOrder    *model.Order
	updatedOrder    *model.Order
	requestedID     int
	deletedID       int
	requestedClient string
}

func (s *serviceStub) Create(order *model.Order) (*model.Order, error) {
	s.createdOrder = order
	order.ID = 1
	order.Status = "new"
	return order, nil
}

func (s *serviceStub) GetByID(id int) (*model.Order, error) {
	s.requestedID = id
	return &model.Order{ID: id}, nil
}

func (s *serviceStub) Update(order *model.Order) (*model.Order, error) {
	s.updatedOrder = order
	return order, nil
}

func (s *serviceStub) Delete(id int) error {
	s.deletedID = id
	return nil
}

func (s *serviceStub) GetItems(id int) ([]model.Item, error) {
	s.requestedID = id
	return []model.Item{{Name: "coffee", Price: 250}}, nil
}

func (s *serviceStub) GetByCustomer(customer string) ([]*model.Order, error) {
	s.requestedClient = customer
	return []*model.Order{{ID: 1, Customer: customer}}, nil
}

func TestCreate(t *testing.T) {
	service := &serviceStub{}
	handler := New(service)
	items := []model.Item{{Name: "coffee", Price: 250}}

	created, err := handler.Create(CreateRequest{Customer: "Alice", Items: items})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if created.ID != 1 || service.createdOrder.Customer != "Alice" {
		t.Fatalf("Create() = %#v, service got %#v", created, service.createdOrder)
	}
	if len(service.createdOrder.Items) != 1 || service.createdOrder.Items[0] != items[0] {
		t.Fatalf("service items = %#v, want %#v", service.createdOrder.Items, items)
	}
}

func TestUpdate(t *testing.T) {
	service := &serviceStub{}
	handler := New(service)

	updated, err := handler.Update(UpdateRequest{
		ID:       2,
		Customer: "Bob",
		Status:   "paid",
		Items:    []model.Item{{Name: "cake", Price: 300}},
	})
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if updated.ID != 2 || service.updatedOrder.Status != "paid" {
		t.Fatalf("Update() = %#v, service got %#v", updated, service.updatedOrder)
	}
}

func TestReadAndDelete(t *testing.T) {
	service := &serviceStub{}
	handler := New(service)

	if _, err := handler.GetByID(3); err != nil || service.requestedID != 3 {
		t.Fatalf("GetByID() error = %v, requested ID = %d", err, service.requestedID)
	}
	if _, err := handler.GetItems(4); err != nil || service.requestedID != 4 {
		t.Fatalf("GetItems() error = %v, requested ID = %d", err, service.requestedID)
	}
	if err := handler.Delete(5); err != nil || service.deletedID != 5 {
		t.Fatalf("Delete() error = %v, deleted ID = %d", err, service.deletedID)
	}
	if _, err := handler.GetByCustomer("Alice"); err != nil || service.requestedClient != "Alice" {
		t.Fatalf("GetByCustomer() error = %v, customer = %q", err, service.requestedClient)
	}
}
