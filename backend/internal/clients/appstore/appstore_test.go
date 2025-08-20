package appstore_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/furqanmk/reviews-browser/config"
	"github.com/furqanmk/reviews-browser/internal/clients/appstore"
	"github.com/furqanmk/reviews-browser/internal/model"
	"github.com/stretchr/testify/require"
)

// mockRoundTripper implements http.RoundTripper for testing.
type mockRoundTripper struct {
	responseBody string
	statusCode   int
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	resp := &http.Response{
		StatusCode: m.statusCode,
		Body:       ioutil.NopCloser(bytes.NewBufferString(m.responseBody)),
		Header:     make(http.Header),
	}
	return resp, nil
}

func TestFetchRecentReviews(t *testing.T) {
	mockJSON := `{
		"feed": {
			"entry": [
				{
					"id": {"label": "1234567890"},
					"author": {"name": {"label": "Apple Sheep 1"}},
					"title": {"label": "Love this app!"},
					"content": {"label": "Great app, works well."},
					"im:rating": {"label": "5"},
					"updated": {"label": "` + time.Now().Format(time.RFC3339) + `"}
				},
				{
					"id": {"label": "1234567891"},
					"author": {"name": {"label": "Apple Sheep 2"}},
					"title": {"label": "Not so great"},
					"content": {"label": "App not working."},
					"im:rating": {"label": "1"},
					"updated": {"label": "` + time.Now().Add(-2*time.Hour).Format(time.RFC3339) + `"}
				},
				{
					"id": {"label": "1234567892"},
					"author": {"name": {"label": "Apple Sheep 3"}},
					"title": {"label": "meh"},
					"content": {"label": "confused"},
					"im:rating": {"label": "3"},
					"updated": {"label": "` + time.Now().Add(-10*time.Hour).Format(time.RFC3339) + `"}
				}
			]
		}
	}`

	mockClient := &http.Client{
		Transport: &mockRoundTripper{
			responseBody: mockJSON,
			statusCode:   http.StatusOK,
		},
	}

	mockCfg := &config.Config{
		AppStoreReviewsURL: "http://example.com/app/%s/page=%d",
		RecencyCutoffHrs:   5,
	}

	c := &appstore.Client{
		HttpClient: mockClient,
		Config:     mockCfg,
	}

	reviews, err := c.FetchRecentReviews("123456")
	if err != nil {
		t.Fatalf("FetchRecentReviews returned error: %v", err)
	}

	if len(reviews) != 2 {
		t.Fatalf("expected 2 reviews, got %d", len(reviews))
	}

	want := []model.Review{
		{
			ID:      "1234567890",
			AppID:   "123456",
			Author:  "Apple Sheep 1",
			Title:   "Love this app!",
			Content: "Great app, works well.",
			Rating:  5,
		},
		{
			ID:      "1234567891",
			AppID:   "123456",
			Author:  "Apple Sheep 2",
			Title:   "Not so great",
			Content: "App not working.",
			Rating:  1,
		},
	}

	for i, r := range reviews {
		require.Equal(t, want[i].ID, r.ID)
		require.Equal(t, want[i].AppID, r.AppID)
		require.Equal(t, want[i].Author, r.Author)
		require.Equal(t, want[i].Title, r.Title)
		require.Equal(t, want[i].Content, r.Content)
		require.Equal(t, want[i].Rating, r.Rating)
	}
}
