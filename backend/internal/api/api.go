package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/furqanmk/reviews-browser/internal/model"
)

type Persistence interface {
	GetRecentReviews(ctx context.Context, appID string) ([]model.Review, error)
}

type API struct {
	db Persistence
}

func NewAPI(db Persistence) *API {
	return &API{db: db}
}

func (a *API) ReadyHandler(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("API is ready")); err != nil {
		log.Println("Write error:", err)
	}
}

// Handler for fetching reviews from the past 48 hours
func (a *API) ReviewsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// extract app id from query parameters
	appID := r.URL.Query().Get("app_id")

	if appID == "" {
		http.Error(w, "Missing app_id", http.StatusBadRequest)
		return
	}

	reviews, err := a.db.GetRecentReviews(ctx, appID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Println("Database error:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(reviews); err != nil {
		http.Error(w, "Encoding error", http.StatusInternalServerError)
		log.Println("Encoding error:", err)
		return
	}
}

// RegisterHandlers registers API endpoints.
func (a *API) RegisterHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/api/ready", a.ReadyHandler)
	mux.HandleFunc("/api/reviews", a.ReviewsHandler)
}
