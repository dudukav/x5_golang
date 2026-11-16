package memory

import (
	"errors"
	"hw5/internal/model"
	"sync"
)

var (
	ErrAlreadyExists = errors.New("order already exists")
	ErrNotFound      = errors.New("order not found")
)

type Repository struct {
	mu     sync.RWMutex
	orders map[int]*model.Order
}

func New() Repository {
	return Repository{
		orders: make(map[int]*model.Order),
	}
}

func (r *Repository) Create(order *model.Order) (*model.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.orders[order.ID]; ok {
		return order, ErrAlreadyExists
	}

	r.orders[order.ID] = order
	return order, nil
}

func (r *Repository) GetByID(id int) (*model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, ok := r.orders[id]
	if !ok {
		return nil, ErrNotFound
	}

	return order, nil
}

func (r *Repository) Update(order *model.Order) (*model.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.orders[order.ID]; !ok {
		return nil, ErrNotFound
	}

	r.orders[order.ID] = order
	return order, nil
}

func (r *Repository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.orders[id]; !ok {
		return ErrNotFound
	}

	delete(r.orders, id)
	return nil
}

func (r *Repository) GetByCustomer(customer string) ([]*model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	orders := make([]*model.Order, 0)
	for _, order := range r.orders {
		if order.Customer == customer {
			orders = append(orders, order)
		}
	}

	if len(orders) == 0 {
		return nil, ErrNotFound
	}

	return orders, nil
}

func (r *Repository) GetItems(id int) ([]model.Item, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, ok := r.orders[id]
	if !ok {
		return nil, ErrNotFound
	}

	items := make([]model.Item, len(order.Items))
	copy(items, order.Items)
	return items, nil
}
