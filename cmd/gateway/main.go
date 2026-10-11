package main

import (
	"log"
	"net/http"
	"github.com/takuya-q/gogate/internal/proxy"
	"github.com/takuya-q/gogate/internal/config"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	gw := proxy.NewGateway(cfg)
	log.Printf("Starting API Gateway on %s", cfg.ListenAddr)
	if err := http.ListenAndServe(cfg.ListenAddr, gw); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}