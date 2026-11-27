package memory

import (
	"sort"
	"sync"

	"hw5/internal/model"
)

var (
	ErrAlreadyExists = model.ErrAlreadyExists
	ErrNotFound      = model.ErrNotFound
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

	stored := cloneOrder(order)
	r.orders[order.ID] = stored
	return cloneOrder(stored), nil
}

func (r *Repository) GetByID(id int) (*model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, ok := r.orders[id]
	if !ok {
		return nil, ErrNotFound
	}

	return cloneOrder(order), nil
}

func (r *Repository) List() ([]*model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	orders := make([]*model.Order, 0, len(r.orders))
	for _, order := range r.orders {
		orders = append(orders, cloneOrder(order))
	}
	sort.Slice(orders, func(i, j int) bool { return orders[i].ID < orders[j].ID })
	return orders, nil
}

func (r *Repository) Update(order *model.Order) (*model.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.orders[order.ID]; !ok {
		return nil, ErrNotFound
	}

	stored := cloneOrder(order)
	r.orders[order.ID] = stored
	return cloneOrder(stored), nil
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
			orders = append(orders, cloneOrder(order))
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

func cloneOrder(order *model.Order) *model.Order {
	if order == nil {
		return nil
	}

	cloned := *order
	cloned.Items = append([]model.Item(nil), order.Items...)
	cloned.History = append([]model.StatusChange(nil), order.History...)
	return &cloned
}
