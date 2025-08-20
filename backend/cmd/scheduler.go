package cmd

import (
	"context"
	"log"

	"github.com/furqanmk/reviews-browser/config"
	"github.com/furqanmk/reviews-browser/internal/cleanup"
	"github.com/furqanmk/reviews-browser/internal/clients/appstore"
	"github.com/furqanmk/reviews-browser/internal/database"
	"github.com/furqanmk/reviews-browser/internal/polling"
)

func StartSchedulers() {
	ctx := context.Background()

	// load up environment variables
	err := config.LoadEnv()
	if err != nil {
		log.Fatalf("failed to load env: %v", err)
	}

	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatalf("failed to get config: %v", err)
	}

	// Initialize database connection
	db, err := database.NewDBConnection(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Initialize App Store client
	appClient := appstore.NewClient(cfg)

	// Start the polling scheduler
	polling := polling.NewPollingScheduler(db, appClient)
	err = polling.Start(ctx)
	if err != nil {
		log.Fatalf("failed to start polling scheduler: %v", err)
	}

	// Start the cleanup scheduler
	cleanup := cleanup.NewCleanupScheduler(db, cfg)
	cleanup.Start(ctx)

	// wait for interrupt signal
	<-ctx.Done()

	// stop the scheduler on done
	defer polling.Stop()
}
