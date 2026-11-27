package main

import (
	"errors"
	"log"
	"log/slog"
	"net/http"
	"time"

	"hw5/internal/config"
	"hw5/internal/service"
	"hw5/internal/storage/memory"
	"hw5/internal/transport"
)

func main() {
	cfg := config.Load()

	repository := memory.New()
	orders := service.New(&repository)
	handler := transport.New(orders)

	server := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           handler.Routes(slog.Default()),
		ReadHeaderTimeout: 5 * time.Second,
	}

	slog.Info("orders service started", "address", cfg.HTTPAddr)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("orders service: %v", err)
	}
}
