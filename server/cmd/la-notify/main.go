// Package main is the entry point for the LaNotify server.
package main

import (
	"fmt"
	"log"
	"net"
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

	serverURL := getServerURL(cfg.Port)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.HomeHandler(serverURL, cfg.Secret))
	mux.HandleFunc("/notification", handler.AuthMiddleware(cfg.Secret, handler.Notification))

	server := &http.Server{
		Addr:              cfg.Port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       time.Duration(cfg.ReadTimeout) * time.Second,
		WriteTimeout:      time.Duration(cfg.WriteTimeout) * time.Second,
		IdleTimeout:       time.Duration(cfg.IdleTimeout) * time.Second,
	}

	log.Printf("LaNotify server started on %s", serverURL)
	log.Printf("Open %s in your browser to see the QR code", serverURL)

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func getServerURL(port string) string {
	ip := getLocalIP()
	if port[0] == ':' {
		return fmt.Sprintf("http://%s%s", ip, port)
	}
	return fmt.Sprintf("http://%s:%s", ip, port)
}

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "localhost"
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "localhost"
}
