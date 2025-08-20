package api_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/furqanmk/reviews-browser/internal/api"
	"github.com/furqanmk/reviews-browser/internal/model"
	"github.com/stretchr/testify/require"
)

// mockPersistence implements the Persistence interface for testing
type mockPersistence struct {
	reviews []model.Review
	err     error
}

func (m *mockPersistence) GetRecentReviews(ctx context.Context, appID string) ([]model.Review, error) {
	return m.reviews, m.err
}

func TestReviewsHandler_Success(t *testing.T) {
	mockReviews := []model.Review{
		{ID: "1", Content: "Great app!"},
		{ID: "2", Content: "Needs improvement."},
	}
	mockAPI := api.NewAPI(&mockPersistence{reviews: mockReviews})

	req := httptest.NewRequest(http.MethodGet, "/api/reviews?app_id=1234", nil)
	w := httptest.NewRecorder()

	mockAPI.ReviewsHandler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	var actual []model.Review
	if err := json.NewDecoder(resp.Body).Decode(&actual); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Len(t, actual, len(mockReviews))

	for i, review := range actual {
		require.Equal(t, mockReviews[i].ID, review.ID)
		require.Equal(t, mockReviews[i].Content, review.Content)
	}
}

func TestReviewsHandler_MissingAppID(t *testing.T) {
	mockAPI := api.NewAPI(&mockPersistence{})

	req := httptest.NewRequest(http.MethodGet, "/api/reviews", nil)
	w := httptest.NewRecorder()

	mockAPI.ReviewsHandler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestReviewsHandler_DBError(t *testing.T) {
	mockAPI := api.NewAPI(&mockPersistence{err: context.DeadlineExceeded})

	req := httptest.NewRequest(http.MethodGet, "/api/reviews?app_id=test-app", nil)
	w := httptest.NewRecorder()

	mockAPI.ReviewsHandler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}
