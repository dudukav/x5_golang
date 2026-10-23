package hw2

import "errors"

var (
	ErrInvalidID    = errors.New("invalid ID")
	ErrNilOrder     = errors.New("order is nil")
	ErrInvalidName  = errors.New("invalid item name")
	ErrInvalidPrice = errors.New("invalid item price")
)

var StatusNames = make(map[string]string)

func init() {
	StatusNames["new"] = "Новый"
	StatusNames["paid"] = "Оплачен"
	StatusNames["cancelled"] = "Отменён"
}

type Order struct {
	ID     int
	Status string
	Items  []Item
}

type Item struct {
	Name  string
	Price float64
}

func NewOrder(id int) (*Order, error) {
	if err := validateID(id); err != nil {
		return nil, err
	}

	return &Order{
		ID:     id,
		Status: StatusNames["new"],
		Items:  make([]Item, 0),
	}, nil
}

func (o *Order) AddItem(item Item) error {
	if o == nil {
		return ErrNilOrder
	}

	if err := validateItem(item); err != nil {
		return err
	}

	o.Items = append(o.Items, item)

	return nil
}

func (o *Order) Total() float64 {
	if len(o.Items) == 0 {
		return 0
	}

	var total float64
	for _, item := range o.Items {
		total += item.Price
	}

	return total
}

func validateID(id int) error {
	if id <= 0 {
		return ErrInvalidID
	}

	return nil
}

func validateItem(item Item) error {
	if item.Name == "" {
		return ErrInvalidName
	}

	if item.Price < 0 {
		return ErrInvalidPrice
	}

	return nil
}
