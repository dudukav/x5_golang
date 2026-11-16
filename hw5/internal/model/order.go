package model

var StatusNames = make(map[string]string)

func init() {
	StatusNames["new"] = "Новый"
	StatusNames["paid"] = "Оплачен"
	StatusNames["cancelled"] = "Отменён"
}

type Order struct {
	ID       int
	Customer string
	Status   string
	Items    []Item
}

type Item struct {
	Name  string
	Price float64
}
