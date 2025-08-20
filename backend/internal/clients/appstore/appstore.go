package appstore

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/furqanmk/reviews-browser/config"
	"github.com/furqanmk/reviews-browser/internal/model"
)

// ReviewResponse models the relevant parts of the App Store reviews API response.
type ReviewResponse struct {
	Feed struct {
		Entry []struct {
			ID     struct{ Label string } `json:"id"`
			Author struct {
				Name struct{ Label string } `json:"name"`
			} `json:"author"`
			Title    struct{ Label string } `json:"title"`
			Content  struct{ Label string } `json:"content"`
			IMRating struct{ Label string } `json:"im:rating"`
			Updated  struct{ Label string } `json:"updated"`
		} `json:"entry"`
	} `json:"feed"`
}

type Client struct {
	HttpClient *http.Client
	Config     *config.Config
}

// NewClient creates a new App Store client.
func NewClient(cfg *config.Config) *Client {
	return &Client{
		HttpClient: &http.Client{Timeout: 10 * time.Second},
		Config:     cfg,
	}
}

// FetchRecentReviews fetches reviews for the given appID, returning only those from the last 48 hours.
func (c *Client) FetchRecentReviews(appID string) ([]model.Review, error) {
	// Backoff strategy constants
	const (
		maxFailures  = 3
		initialDelay = 5 * time.Second
	)

	var (
		reviews       []model.Review
		page          = 1
		now           = time.Now()
		recencyCutOff = time.Duration(c.Config.RecencyCutoffHrs) * time.Hour
	)

	for {
		var resp ReviewResponse
		var body []byte
		var err error

		// Rudimentary exponential backoff
		for attempt := 0; attempt < maxFailures; attempt++ {
			url := fmt.Sprintf(c.Config.AppStoreReviewsURL, appID, page)
			httpResp, reqErr := c.HttpClient.Get(url)
			if reqErr != nil {
				delay := initialDelay * (time.Duration(attempt + 1))
				time.Sleep(delay)
				continue
			}
			defer httpResp.Body.Close()
			body, err = io.ReadAll(httpResp.Body)
			if err != nil {
				delay := initialDelay * (time.Duration(attempt + 1))
				time.Sleep(delay)
				continue
			}
			if httpResp.StatusCode != http.StatusOK {
				delay := initialDelay * (time.Duration(attempt + 1))
				time.Sleep(delay)
				continue
			}
			break
		}
		if err != nil {
			return nil, errors.New("failed to fetch reviews after retries")
		}

		if err := json.Unmarshal(body, &resp); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}

		entries := resp.Feed.Entry
		if len(entries) == 0 {
			break
		}

		stop := false
		for _, entry := range entries {
			updated, err := time.Parse(time.RFC3339, entry.Updated.Label)
			if err != nil {
				continue
			}
			if now.Sub(updated) > recencyCutOff {
				stop = true
				break
			}
			rating := 0
			if _, err := fmt.Sscanf(entry.IMRating.Label, "%d", &rating); err != nil {
				continue
			}
			review := model.Review{
				ID:        entry.ID.Label,
				AppID:     appID,
				Author:    entry.Author.Name.Label,
				Title:     entry.Title.Label,
				Content:   entry.Content.Label,
				Rating:    rating,
				CreatedAt: updated,
			}
			reviews = append(reviews, review)
		}
		if stop {
			break
		}
		page++
	}

	return reviews, nil
}
