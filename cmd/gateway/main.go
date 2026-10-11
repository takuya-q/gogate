package main

import (
	"log"

	"github.com/takuya-q/gogate/internal/config"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Printf("Starting API Gateway on %s", cfg.ListenAddr)
	log.Printf("Loaded %d routes", len(cfg.Routes))
	
	for _, route := range cfg.Routes {
		log.Printf("Route: %s -> %d backends", route.Path, len(route.Backends))
	}
}