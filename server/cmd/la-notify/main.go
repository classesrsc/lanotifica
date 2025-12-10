// Package main is the entry point for the LA-notify server.
package main

import (
	"log"
	"net/http"
	"time"

	"github.com/alessandrolattao/la-notify/internal/config"
	"github.com/alessandrolattao/la-notify/internal/handler"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Printf("Config loaded from %s", config.Path())

	mux := http.NewServeMux()
	mux.HandleFunc("/notification", handler.Notification)

	server := &http.Server{
		Addr:              cfg.Port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       time.Duration(cfg.ReadTimeout) * time.Second,
		WriteTimeout:      time.Duration(cfg.WriteTimeout) * time.Second,
		IdleTimeout:       time.Duration(cfg.IdleTimeout) * time.Second,
	}

	log.Printf("LA-notify server started on http://localhost%s", server.Addr)

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
