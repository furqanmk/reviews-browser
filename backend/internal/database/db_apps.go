package database

import (
	"context"
	"strconv"
	"time"

	"github.com/furqanmk/reviews-browser/internal/model"
)

const (
	COLUMN_APPS_ID = iota
	COLUMN_APPS_LAST_POLLED
	COLUMN_APPS_POLL_EVERY
)

var (
	appsHeader = []string{
		"id",
		"last_fetched",
		"poll_every_seconds",
	}
)

// GetApps retrieves all apps from the apps CSV file.
func (db *DB) GetApps(ctx context.Context) ([]model.App, error) {
	reader, file, err := getReader(db.config.AppsCSV)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var apps []model.App
	for _, row := range rows[1:] {
		lastFetched, err := time.Parse(time.RFC3339, row[COLUMN_APPS_LAST_POLLED])
		if err != nil {
			continue
		}
		pollEverySeconds, err := strconv.Atoi(row[COLUMN_APPS_POLL_EVERY])
		if err != nil {
			continue
		}
		apps = append(apps, model.App{
			ID:               row[COLUMN_APPS_ID],
			LastFetched:      lastFetched,
			PollEverySeconds: pollEverySeconds,
		})
	}

	return apps, nil
}

func (db *DB) UpdateLastFetched(ctx context.Context, appID string, lastFetched time.Time) error {
	// Update the last fetched time for the app in the database
	apps, err := db.GetApps(ctx)
	if err != nil {
		return err
	}

	for i, app := range apps {
		if app.ID == appID {
			apps[i].LastFetched = lastFetched
			break
		}
	}

	writer, file, err := emptyFile(db.config.AppsCSV, appsHeader)
	if err != nil {
		return err
	}
	defer file.Close()
	defer writer.Flush()

	for _, app := range apps {
		record := []string{
			app.ID,
			app.LastFetched.Format(time.RFC3339),
			strconv.Itoa(app.PollEverySeconds),
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}
	return nil
}
