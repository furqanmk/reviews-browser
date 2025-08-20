package polling

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/furqanmk/reviews-browser/internal/clients/appstore"
	"github.com/furqanmk/reviews-browser/internal/database"
	"github.com/furqanmk/reviews-browser/internal/model"
)

// PollingScheduler manages polling of app reviews.
type PollingScheduler struct {
	db        *database.DB
	appClient *appstore.Client
	wg        sync.WaitGroup
	stopCh    chan struct{}
}

// NewPollingAgent creates a new polling agent.
func NewPollingScheduler(DB *database.DB, AppClient *appstore.Client) *PollingScheduler {
	return &PollingScheduler{
		db:        DB,
		appClient: AppClient,
		stopCh:    make(chan struct{}),
	}
}

// Start begins polling for all apps.
func (a *PollingScheduler) Start(ctx context.Context) error {
	log.Print("Starting polling scheduler...")

	apps, err := a.db.GetApps(ctx)
	if err != nil {
		return err
	}

	for _, app := range apps {
		a.wg.Add(1)
		go a.pollForReviews(ctx, app)
	}
	return nil
}

// pollForReviews schedules polling for a single app.
func (a *PollingScheduler) pollForReviews(ctx context.Context, app model.App) {
	lastFetched := app.LastFetched

	defer a.wg.Done()
	for {
		nextPoll := lastFetched.Add(time.Duration(app.PollEverySeconds) * time.Second)
		wait := time.Until(nextPoll)
		if wait > 0 {
			select {
			case <-time.After(wait):
			case <-a.stopCh:
				break
			}
		}

		// Fetch recent reviews and update data store
		reviews, err := a.appClient.FetchRecentReviews(app.ID)
		if err != nil {
			log.Printf("Error fetching reviews for app %s: %v", app.ID, err)
		} else {
			for _, review := range reviews {
				_ = a.db.InsertReview(ctx, review, app.ID)
			}
		}

		// Update last fetched time
		lastFetched = time.Now()
		err = a.db.UpdateLastFetched(ctx, app.ID, lastFetched)
		if err != nil {
			log.Printf("Error updating last fetched time for app %s: %v", app.ID, err)
		}
	}
}

// Stop halts all polling goroutines.
func (a *PollingScheduler) Stop() {
	close(a.stopCh)
	a.wg.Wait()
}
