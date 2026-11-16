package transport

import "hw5/internal/model"

type Handler struct {
	orders OrderService
}

type CreateRequest struct {
	Customer string
	Items    []model.Item
}

type UpdateRequest struct {
	ID       int
	Customer string
	Status   string
	Items    []model.Item
}

func New(orders OrderService) *Handler {
	return &Handler{orders: orders}
}

func (h *Handler) Create(request CreateRequest) (*model.Order, error) {
	return h.orders.Create(&model.Order{
		Customer: request.Customer,
		Items:    request.Items,
	})
}

func (h *Handler) GetByID(id int) (*model.Order, error) {
	return h.orders.GetByID(id)
}

func (h *Handler) Update(request UpdateRequest) (*model.Order, error) {
	return h.orders.Update(&model.Order{
		ID:       request.ID,
		Customer: request.Customer,
		Status:   request.Status,
		Items:    request.Items,
	})
}

func (h *Handler) Delete(id int) error {
	return h.orders.Delete(id)
}

func (h *Handler) GetItems(id int) ([]model.Item, error) {
	return h.orders.GetItems(id)
}

func (h *Handler) GetByCustomer(customer string) ([]*model.Order, error) {
	return h.orders.GetByCustomer(customer)
}
