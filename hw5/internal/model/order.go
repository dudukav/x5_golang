package model

import "time"

type Status string

const (
	StatusNew        Status = "new"
	StatusProcessing Status = "processing"
	StatusCompleted  Status = "completed"
	StatusCancelled  Status = "cancelled"
)

var StatusNames = map[Status]string{
	StatusNew:        "Новый",
	StatusProcessing: "В обработке",
	StatusCompleted:  "Завершён",
	StatusCancelled:  "Отменён",
}

type Order struct {
	ID         int            `json:"id"`
	Customer   string         `json:"customer,omitempty"`
	Items      []Item         `json:"items"`
	TotalPrice int64          `json:"total_price"`
	Status     Status         `json:"status"`
	CreatedAt  time.Time      `json:"created_at"`
	History    []StatusChange `json:"history,omitempty"`
}

type Item struct {
	Name     string `json:"name"`
	Quantity int64  `json:"quantity"`
	Price    int64  `json:"price"`
}

type StatusChange struct {
	Status    Status    `json:"status"`
	ChangedAt time.Time `json:"changed_at"`
}

func IsValidStatus(status Status) bool {
	_, ok := StatusNames[status]
	return ok
}
