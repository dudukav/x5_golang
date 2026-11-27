package service

import "hw5/internal/model"

type Repository interface {
	Create(*model.Order) (*model.Order, error)
	GetByID(id int) (*model.Order, error)
	List() ([]*model.Order, error)
	Update(order *model.Order) (*model.Order, error)
	Delete(id int) error
	GetItems(id int) ([]model.Item, error)
	GetByCustomer(customer string) ([]*model.Order, error)
}
