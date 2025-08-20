package cleanup

import (
	"context"
	"log"
	"time"

	"github.com/furqanmk/reviews-browser/config"
)

// CleanupScheduler manages cleanup of old reviews.
type CleanupScheduler struct {
	db  CleanupDB
	cfg *config.Config
}

type CleanupDB interface {
	CleanUpOldReviews(ctx context.Context) error
}

func NewCleanupScheduler(db CleanupDB, cfg *config.Config) *CleanupScheduler {
	return &CleanupScheduler{
		db:  db,
		cfg: cfg,
	}
}

// StartCleanupScheduler runs CleanUpOldReviews every number of hours configured in the env.
func (s *CleanupScheduler) Start(ctx context.Context) {
	log.Print("Starting cleanup scheduler...")

	go func() {
		for {
			err := s.db.CleanUpOldReviews(ctx)
			if err != nil {
				log.Printf("Error cleaning up old reviews: %v", err)
			}
			time.Sleep(time.Duration(s.cfg.CleanupEveryHrs) * time.Hour)
		}
	}()
}
