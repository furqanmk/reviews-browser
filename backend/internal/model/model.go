package model

import "time"

// App represents an application record.
type App struct {
	ID               string
	LastFetched      time.Time
	PollEverySeconds int
}

// Review represents a review record.
type Review struct {
	ID        string    `json:"id"`
	AppID     string    `json:"app_id"`
	Author    string    `json:"author"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Rating    int       `json:"rating"`
	CreatedAt time.Time `json:"created_at"`
}
