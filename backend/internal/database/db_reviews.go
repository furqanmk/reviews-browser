package database

import (
	"context"
	"strconv"
	"time"

	"github.com/furqanmk/reviews-browser/internal/model"
)

const (
	COLUMN_REVIEWS_ID = iota
	COLUMN_REVIEWS_APP_ID
	COLUMN_REVIEWS_AUTHOR
	COLUMN_REVIEWS_TITLE
	COLUMN_REVIEWS_CONTENT
	COLUMN_REVIEWS_RATING
	COLUMN_REVIEWS_DATE
)

var (
	reviewsHeader = []string{
		"id",
		"app_id",
		"author",
		"title",
		"content",
		"rating",
		"date",
	}
)

// GetRecentReviews retrieves reviews for a specific app ID from the last x hours.
func (db *DB) GetRecentReviews(ctx context.Context, appID string) ([]model.Review, error) {
	reader, file, err := getReader(db.config.ReviewsCSV)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var reviews []model.Review
	cutoff := time.Now().Add(-time.Duration(db.config.RecencyCutoffHrs) * time.Hour)

	for _, row := range rows[1:] {
		// Ensure row has enough columns
		if len(row) < 7 {
			continue
		}

		// Filter out reviews by App ID
		if row[COLUMN_REVIEWS_APP_ID] != appID {
			continue
		}

		// Parse the review creation date
		createdAt, err := time.Parse(time.RFC3339, row[COLUMN_REVIEWS_DATE])
		if err != nil || createdAt.Before(cutoff) {
			continue
		}

		// Parse the review rating
		rating, err := strconv.Atoi(row[COLUMN_REVIEWS_RATING])
		if err != nil {
			continue
		}

		reviews = append(reviews, model.Review{
			ID:        row[COLUMN_REVIEWS_ID],
			AppID:     row[COLUMN_REVIEWS_APP_ID],
			Author:    row[COLUMN_REVIEWS_AUTHOR],
			Title:     row[COLUMN_REVIEWS_TITLE],
			Content:   row[COLUMN_REVIEWS_CONTENT],
			Rating:    rating,
			CreatedAt: createdAt,
		})
	}

	// Sort reviews by CreatedAt descending
	for i := range len(reviews) - 1 {
		for j := i + 1; j < len(reviews); j++ {
			if reviews[i].CreatedAt.Before(reviews[j].CreatedAt) {
				reviews[i], reviews[j] = reviews[j], reviews[i]
			}
		}
	}

	return reviews, nil
}

// InsertReview adds a new review to the reviews CSV file and updates the in-memory data.
func (db *DB) InsertReview(ctx context.Context, review model.Review, appID string) error {
	existingReviews, err := db.GetRecentReviews(ctx, appID)
	if err != nil {
		return err
	}

	// skip writing if the review already exists
	for _, r := range existingReviews {
		if r.ID == review.ID {
			return nil
		}
	}

	// Write the new review to the end of the CSV file
	writer, file, err := getWriter(db.config.ReviewsCSV)
	if err != nil {
		return err
	}
	defer file.Close()
	defer writer.Flush()

	record := []string{
		review.ID,
		appID,
		review.Author,
		review.Title,
		review.Content,
		strconv.Itoa(review.Rating),
		review.CreatedAt.Format(time.RFC3339),
	}

	if err := writer.Write(record); err != nil {
		return err
	}

	return nil
}

func (db *DB) CleanUpOldReviews(ctx context.Context) error {
	// Read all reviews from CSV
	reader, file, err := getReader(db.config.ReviewsCSV)
	if err != nil {
		return err
	}
	defer file.Close()

	rows, err := reader.ReadAll()
	if err != nil {
		return err
	}

	// Prepare cutoff time
	cutoff := time.Now().Add(-time.Duration(db.config.RecencyCutoffHrs) * time.Hour)

	// Filter out the older reviews
	var newRows [][]string
	for _, row := range rows[1:] {
		if len(row) < 7 {
			continue
		}
		createdAt, err := time.Parse(time.RFC3339, row[COLUMN_REVIEWS_DATE])
		if err != nil {
			continue
		}
		if !createdAt.Before(cutoff) {
			newRows = append(newRows, row)
		}
	}

	writer, file, err := emptyFile(db.config.ReviewsCSV, reviewsHeader)
	if err != nil {
		return err
	}
	defer file.Close()
	defer writer.Flush()

	for _, row := range newRows {
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}
