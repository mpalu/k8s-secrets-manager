package main

import (
	"log"
	"os"

	"github.com/mpalu/k8s-secrets-manager/internal/cli/cmd"
	"github.com/mpalu/k8s-secrets-manager/internal/config"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Execute command with configuration
	if err := cmd.Execute(cfg); err != nil {
		log.Printf("Error executing command: %v", err)
		os.Exit(1)
	}
}
