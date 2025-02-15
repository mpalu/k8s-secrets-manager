package main

import (
	"log"
	"os"

	"github.com/mpalu/k8s-secrets-manager/internal/cli/cmd"
	"github.com/mpalu/k8s-secrets-manager/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	if err := cmd.Execute(cfg); err != nil {
		os.Exit(1)
	}
}
