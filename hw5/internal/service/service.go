package service

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"sync"
	"time"

	"hw5/internal/model"
)

var (
	ErrInvalidOrder            = errors.New("invalid order")
	ErrInvalidStatusTransition = errors.New("status transition is not allowed")
)

type Service struct {
	repository Repository
	mu         sync.Mutex
	id         int
	now        func() time.Time
}

func New(repository Repository) *Service {
	return &Service{
		repository: repository,
		id:         1,
		now:        time.Now,
	}
}

func (s *Service) Create(order *model.Order) (*model.Order, error) {
	if err := validateOrder(order, false); err != nil {
		return nil, err
	}

	totalPrice, err := calculateTotal(order.Items)
	if err != nil {
		return nil, err
	}

	now := s.now().UTC()
	orderToCreate := cloneOrder(order)
	orderToCreate.ID = s.nextID()
	orderToCreate.Status = model.StatusNew
	orderToCreate.TotalPrice = totalPrice
	orderToCreate.CreatedAt = now
	orderToCreate.History = []model.StatusChange{{Status: model.StatusNew, ChangedAt: now}}

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

func (s *Service) List() ([]*model.Order, error) {
	orders, err := s.repository.List()
	if err != nil {
		return nil, fmt.Errorf("list orders: %w", err)
	}
	return orders, nil
}

func (s *Service) UpdateStatus(id int, status model.Status) (*model.Order, error) {
	if err := validateID(id); err != nil {
		return nil, err
	}
	if !model.IsValidStatus(status) {
		return nil, fmt.Errorf("%w: unknown status %q", ErrInvalidOrder, status)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	order, err := s.repository.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("get order %d before status update: %w", id, err)
	}
	if !canTransition(order.Status, status) {
		return nil, fmt.Errorf("%w: %s -> %s", ErrInvalidStatusTransition, order.Status, status)
	}

	order.Status = status
	order.History = append(order.History, model.StatusChange{
		Status:    status,
		ChangedAt: s.now().UTC(),
	})

	updated, err := s.repository.Update(order)
	if err != nil {
		return nil, fmt.Errorf("update order %d status: %w", id, err)
	}
	return updated, nil
}

func (s *Service) History(id int) ([]model.StatusChange, error) {
	order, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}
	return append([]model.StatusChange(nil), order.History...), nil
}

func (s *Service) Update(order *model.Order) (*model.Order, error) {
	if err := validateOrder(order, true); err != nil {
		return nil, err
	}

	totalPrice, err := calculateTotal(order.Items)
	if err != nil {
		return nil, err
	}
	orderToUpdate := cloneOrder(order)
	orderToUpdate.TotalPrice = totalPrice

	updated, err := s.repository.Update(orderToUpdate)
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
		if !model.IsValidStatus(order.Status) {
			return fmt.Errorf("%w: unknown status %q", ErrInvalidOrder, order.Status)
		}
	}
	if len(order.Items) == 0 {
		return fmt.Errorf("%w: items are empty", ErrInvalidOrder)
	}
	for i, item := range order.Items {
		if strings.TrimSpace(item.Name) == "" {
			return fmt.Errorf("%w: item %d name is empty", ErrInvalidOrder, i)
		}
		if item.Quantity <= 0 {
			return fmt.Errorf("%w: item %d quantity must be positive", ErrInvalidOrder, i)
		}
		if item.Price <= 0 {
			return fmt.Errorf("%w: item %d price must be positive", ErrInvalidOrder, i)
		}
	}
	return nil
}

func calculateTotal(items []model.Item) (int64, error) {
	var total int64
	for i, item := range items {
		if item.Price > math.MaxInt64/item.Quantity || total > math.MaxInt64-item.Price*item.Quantity {
			return 0, fmt.Errorf("%w: item %d makes total price overflow", ErrInvalidOrder, i)
		}
		total += item.Price * item.Quantity
	}
	return total, nil
}

func canTransition(from, to model.Status) bool {
	switch from {
	case model.StatusNew:
		return to == model.StatusProcessing || to == model.StatusCancelled
	case model.StatusProcessing:
		return to == model.StatusCompleted || to == model.StatusCancelled
	default:
		return false
	}
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
	cloned.History = append([]model.StatusChange(nil), order.History...)
	return &cloned
}
