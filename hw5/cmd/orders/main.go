package main

import (
	"fmt"
	"log"

	"hw5/internal/config"
	"hw5/internal/model"
	"hw5/internal/service"
	"hw5/internal/storage/memory"
	"hw5/internal/transport"
)

func main() {
	cfg := config.Load()

	repository := memory.New()
	orders := service.New(&repository)
	handler := transport.New(orders)

	order, err := handler.Create(transport.CreateRequest{
		Customer: "Анна",
		Items: []model.Item{
			{Name: "Кофе", Price: 250},
		},
	})
	if err != nil {
		log.Fatalf("создание заказа: %v", err)
	}

	fmt.Printf("сервис слушает на порту: %s, создан заказ: %+v\n", cfg.HTTPAddr, order)
}
