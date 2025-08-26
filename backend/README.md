# Reviews Browser Go Backend

## Overview
Backend service for the reviews browser app that lets users browse recent App Store reviews for iOS apps. The service contains two distinct components: 
1. a polling service that periodically fetches App Store reviews for specified iOS apps
2. an API service that provides a REST API to access recent reviews.

## Service Architecture

### 1. Scheduler Component
#### Polling Scheduler
- Fetches reviews from App Store Connect RSS feed
- Stores new reviews in CSV file that acts as the reviews table
- Maintains last polled timestamp

#### Cleanup Scheduler
- Purges reviews older than 48 hours (configurable through environment variable)

### 2. API Service Component
- Returns filtered reviews based on the app ID in the query parameters

**Endpoint**: `GET /api/reviews?app_id={appId}`

**Response**:
```json
[
  {
    "id": "12345",
    "author": "User1",
    "rating": 5,
    "title": "Great app!",
    "content": "Works perfectly",
    "date": "2023-11-15T12:00:00Z"
  },
  // ...
]
```

## Deployment Considerations

1. **Single Instance**: Run as a single process (no clustering needed)
2. **Persistent Storage**: Ensure CSV files have proper file permissions
3. **Logging**: Implement basic logging for polling activities
4. **Error Handling**: Retry logic for failed RSS fetches

## Database Schema
With a real database, we could use the following table schema, for this case we have used CSV files to simulate the database tables.

```sql
-- Table to store apps being tracked
CREATE TABLE apps (
    app_id INTEGER PRIMARY KEY,
    last_polled TIMESTAMP NOT NULL,
    poll_every_seconds INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(app_id)
);

-- Table to store reviews
CREATE TABLE reviews (
    id SERIAL PRIMARY KEY,
    app_id INTEGER NOT NULL,
    author TEXT NOT NULL,
    rating SMALLINT NOT NULL,
    title TEXT,
    content TEXT,
    review_date TIMESTAMP NOT NULL
);
```

## Folder Structure

```
app-review-browser/
├── cmd/                # Main application entrypoints
│   ├── api.go          # HTTP server main.go
│   └── scheduler.go    # Polling scheduler
├── config/             # Configuration files
├── data/               # CSV data files
├── docs/               # Documentation files
├── internal/           # Private application code
│   ├── api/            # HTTP handlers and routing
│   ├── cleanup/        # Logic for cleaning up older reviews
│   ├── client/         # HTTP client for fetching reviews
│       └── appstore.go # App Store Connect API client
│   ├── database/       # Database access and models
│   ├── model/          # Data models
│   └── polling/        # Polling logic and RSS fetching
├── go.mod
├── go.sum
└── README.md
```

## Follow Up
- Front-end can query an app ID
- Backend fetches reviews based on an app ID