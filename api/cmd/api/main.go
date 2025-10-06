package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/yeftaz/susano.id/api/internal/config"
	"github.com/yeftaz/susano.id/api/internal/database"
	"github.com/yeftaz/susano.id/api/internal/router"
	appLogger "github.com/yeftaz/susano.id/api/pkg/logger"
)

func main() {
	// Parse command line flags
	showRoutes := flag.Bool("routes", false, "Show all API routes")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	logger := appLogger.New(cfg)
	defer logger.Sync()

	// Connect to database
	db, err := database.Connect(cfg)
	if err != nil {
		logger.Fatal("Failed to connect to database", "error", err)
	}
	defer db.Close()

	logger.Info("Database connection established")

	// Initialize router
	r := router.New(cfg, db, logger)

	// If -routes flag is set, show routes and exit
	if *showRoutes {
		router.ShowRoutes(r)
		os.Exit(0)
	}

	// Start server
	addr := fmt.Sprintf(":%s", cfg.AppPort)
	logger.Info("Starting server", "address", addr, "environment", cfg.AppEnv)

	if err := http.ListenAndServe(addr, r); err != nil {
		logger.Fatal("Server failed to start", "error", err)
	}
}
