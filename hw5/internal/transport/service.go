package transport

import "hw5/internal/model"

type OrderService interface {
	Create(order *model.Order) (*model.Order, error)
	GetByID(id int) (*model.Order, error)
	List() ([]*model.Order, error)
	UpdateStatus(id int, status model.Status) (*model.Order, error)
	History(id int) ([]model.StatusChange, error)
}
