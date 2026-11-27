package transport

import (
	"context"

	"hw5/internal/model"
)

type Server struct {
	orders OrderService
}

var _ StrictServerInterface = (*Server)(nil)

func New(orders OrderService) *Server {
	return &Server{orders: orders}
}

func (h *Server) Health(context.Context, HealthRequestObject) (HealthResponseObject, error) {
	status := "ok"
	return Health200JSONResponse{Status: &status}, nil
}

func (h *Server) CreateOrder(_ context.Context, request CreateOrderRequestObject) (CreateOrderResponseObject, error) {
	items := make([]model.Item, len(request.Body.Items))
	for i, item := range request.Body.Items {
		items[i] = model.Item{Name: item.Name, Quantity: item.Quantity, Price: item.Price}
	}

	order, err := h.orders.Create(&model.Order{
		Customer: stringValue(request.Body.Customer),
		Items:    items,
	})
	if err != nil {
		return nil, err
	}
	return CreateOrder201JSONResponse(orderResponse(order)), nil
}

func (h *Server) ListOrders(context.Context, ListOrdersRequestObject) (ListOrdersResponseObject, error) {
	orders, err := h.orders.List()
	if err != nil {
		return nil, err
	}

	response := make(ListOrders200JSONResponse, len(orders))
	for i, order := range orders {
		response[i] = orderResponse(order)
	}
	return response, nil
}

func (h *Server) GetOrder(_ context.Context, request GetOrderRequestObject) (GetOrderResponseObject, error) {
	order, err := h.orders.GetByID(request.Id)
	if err != nil {
		return nil, err
	}
	return GetOrder200JSONResponse(orderResponse(order)), nil
}

func (h *Server) UpdateOrderStatus(_ context.Context, request UpdateOrderStatusRequestObject) (UpdateOrderStatusResponseObject, error) {
	order, err := h.orders.UpdateStatus(request.Id, model.Status(request.Body.Status))
	if err != nil {
		return nil, err
	}
	return UpdateOrderStatus200JSONResponse(orderResponse(order)), nil
}

func (h *Server) GetOrderHistory(_ context.Context, request GetOrderHistoryRequestObject) (GetOrderHistoryResponseObject, error) {
	history, err := h.orders.History(request.Id)
	if err != nil {
		return nil, err
	}

	response := make(GetOrderHistory200JSONResponse, len(history))
	for i, change := range history {
		response[i] = statusChangeResponse(change)
	}
	return response, nil
}

func orderResponse(order *model.Order) Order {
	items := make([]Item, len(order.Items))
	for i, item := range order.Items {
		items[i] = Item{Name: item.Name, Quantity: item.Quantity, Price: item.Price}
	}
	history := make([]StatusChange, len(order.History))
	for i, change := range order.History {
		history[i] = statusChangeResponse(change)
	}

	var customer *string
	if order.Customer != "" {
		customer = &order.Customer
	}
	return Order{
		Id:         order.ID,
		Customer:   customer,
		Items:      items,
		TotalPrice: order.TotalPrice,
		Status:     Status(order.Status),
		CreatedAt:  order.CreatedAt,
		History:    history,
	}
}

func statusChangeResponse(change model.StatusChange) StatusChange {
	return StatusChange{Status: Status(change.Status), ChangedAt: change.ChangedAt}
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
