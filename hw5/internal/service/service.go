package service

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"hw5/internal/model"
)

var ErrInvalidOrder = errors.New("invalid order")

type Service struct {
	repository Repository
	mu         sync.Mutex
	id         int
}

func New(repository Repository) *Service {
	return &Service{
		repository: repository,
		id:         1,
	}
}

func (s *Service) Create(order *model.Order) (*model.Order, error) {
	if err := validateOrder(order, false); err != nil {
		return nil, err
	}

	orderToCreate := cloneOrder(order)
	orderToCreate.ID = s.nextID()
	orderToCreate.Status = "new"

	created, err := s.repository.Create(orderToCreate)
	if err != nil {
		return nil, fmt.Errorf("create order: %w", err)
	}
	return created, nil
}

func (s *Service) GetByID(id int) (*model.Order, error) {
	if err := validateID(id); err != nil {
		return nil, err
	}

	order, err := s.repository.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("get order %d: %w", id, err)
	}
	return order, nil
}

func (s *Service) Update(order *model.Order) (*model.Order, error) {
	if err := validateOrder(order, true); err != nil {
		return nil, err
	}

	updated, err := s.repository.Update(cloneOrder(order))

	if err != nil {
		return nil, fmt.Errorf("update order %d: %w", order.ID, err)
	}

	return updated, nil
}

func (s *Service) Delete(id int) error {
	if err := validateID(id); err != nil {
		return err
	}

	if err := s.repository.Delete(id); err != nil {
		return fmt.Errorf("delete order %d: %w", id, err)
	}

	return nil
}

func (s *Service) GetItems(id int) ([]model.Item, error) {
	if err := validateID(id); err != nil {
		return nil, err
	}

	items, err := s.repository.GetItems(id)
	if err != nil {
		return nil, fmt.Errorf("get items for order %d: %w", id, err)
	}

	return items, nil
}

func (s *Service) GetByCustomer(customer string) ([]*model.Order, error) {
	if strings.TrimSpace(customer) == "" {
		return nil, fmt.Errorf("%w: customer is empty", ErrInvalidOrder)
	}

	orders, err := s.repository.GetByCustomer(customer)
	if err != nil {
		return nil, fmt.Errorf("get orders for customer %q: %w", customer, err)
	}

	return orders, nil
}

func (s *Service) nextID() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := s.id
	s.id++
	return id
}

func validateOrder(order *model.Order, requireID bool) error {
	if order == nil {
		return fmt.Errorf("%w: order is nil", ErrInvalidOrder)
	}

	if requireID {
		if err := validateID(order.ID); err != nil {
			return err
		}

		if _, ok := model.StatusNames[order.Status]; !ok {
			return fmt.Errorf("%w: unknown status %q", ErrInvalidOrder, order.Status)
		}
	}

	if strings.TrimSpace(order.Customer) == "" {
		return fmt.Errorf("%w: customer is empty", ErrInvalidOrder)
	}

	if len(order.Items) == 0 {
		return fmt.Errorf("%w: items are empty", ErrInvalidOrder)
	}

	for i, item := range order.Items {
		if strings.TrimSpace(item.Name) == "" {
			return fmt.Errorf("%w: item %d name is empty", ErrInvalidOrder, i)
		}

		if item.Price <= 0 {
			return fmt.Errorf("%w: item %d price must be positive", ErrInvalidOrder, i)
		}
	}

	return nil
}

func validateID(id int) error {
	if id <= 0 {
		return fmt.Errorf("%w: id must be positive", ErrInvalidOrder)
	}

	return nil
}

func cloneOrder(order *model.Order) *model.Order {
	cloned := *order
	cloned.Items = append([]model.Item(nil), order.Items...)
	return &cloned
}
