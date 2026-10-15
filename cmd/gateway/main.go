package main

import (
	"log"
	"net/http"

	"github.com/takuya-q/gogate/internal/config"
	"github.com/takuya-q/gogate/internal/middleware"
	"github.com/takuya-q/gogate/internal/proxy"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	gw := proxy.NewGateway(cfg)
	handler := middleware.RateLimit(gw)
	log.Printf("Starting API Gateway on %s", cfg.ListenAddr)
	if err := http.ListenAndServe(cfg.ListenAddr, handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}