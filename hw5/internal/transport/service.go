package transport

import "hw5/internal/model"

type OrderService interface {
	Create(order *model.Order) (*model.Order, error)
	GetByID(id int) (*model.Order, error)
	Update(order *model.Order) (*model.Order, error)
	Delete(id int) error
	GetItems(id int) ([]model.Item, error)
	GetByCustomer(customer string) ([]*model.Order, error)
}