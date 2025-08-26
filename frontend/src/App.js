import React, { useState, useEffect } from 'react';
import './App.css';

// Star rating component
function StarRating({ rating }) {
    const maxStars = 5;
    const stars = [];
    
    for (let i = 1; i <= maxStars; i++) {
        stars.push(
            <i 
                key={i}
                className={i <= rating ? "bi bi-star-fill star-filled me-1" : "bi bi-star star-empty me-1"}
            ></i>
        );
    }
    
    return <div className="d-flex">{stars}</div>;
}

// Review card component
function ReviewCard({ review }) {
    const formatDate = (dateString) => {
        const date = new Date(dateString);
        return date.toLocaleDateString('en-US', {
            year: 'numeric',
            month: 'short',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
        });
    };

    return (
        <div className="card review-card mb-3">
            <div className="card-body">
                <div className="d-flex justify-content-between align-items-start mb-2">
                    <div className="d-flex align-items-center">
                        <div className="bg-primary text-white rounded-circle d-flex align-items-center justify-content-center me-2" 
                             style={{width: '40px', height: '40px', fontSize: '14px', fontWeight: 'bold'}}>
                            {review.author ? review.author.charAt(0).toUpperCase() : 'U'}
                        </div>
                        <div>
                            <h6 className="card-title mb-0">{review.author || 'Anonymous'}</h6>
                            <small className="text-muted">{formatDate(review.created_at)}</small>
                        </div>
                    </div>
                    <span className="badge bg-primary rating-badge">
                        <i className="bi bi-star-fill me-1"></i>
                        {review.rating}
                    </span>
                </div>
                
                <StarRating rating={review.rating} />
                
                {review.title && (
                    <h6 className="card-subtitle mt-2 mb-1 text-dark">{review.title}</h6>
                )}
                
                <p className="card-text mt-2">{review.content}</p>
            </div>
        </div>
    );
}

// Main app component
function App() {
    const [reviews, setReviews] = useState([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);
    const [appId, setAppId] = useState('');
    const [inputAppId, setInputAppId] = useState('');

    const fetchReviews = async (appID) => {
        if (!appID.trim()) {
            setError('Please enter a valid App ID');
            return;
        }

        setLoading(true);
        setError(null);
        
        try {
            const url = `/api/reviews_by_app?app_id=${appID}`;
            const response = await fetch(url);
            const reviews = await response.json();

            if (reviews == null) {
              setReviews([]);
            } else {
              setReviews(reviews);
            }

            setAppId(appID);
            
        } catch (err) {
            setError('Failed to fetch reviews. Please try again.');
            console.error('Error fetching reviews:', err);
        } finally {
            setLoading(false);
        }
    };

    const handleSubmit = (e) => {
        e.preventDefault();
        fetchReviews(inputAppId.trim());
    };

    const handleQuickAppSelect = (quickAppId) => {
        setInputAppId(quickAppId);
        fetchReviews(quickAppId);
    };

    return (
        <div className="container mt-4">
            <div className="row justify-content-center">
                <div className="col-lg-10 col-xl-8">
                    {/* Header */}
                    <div className="text-center mb-4">
                        <h1 className="h2 mb-3">
                            <i className="bi bi-chat-square-text-fill text-primary me-2"></i>
                            App Store Reviews
                        </h1>
                        
                        {/* App ID Input Form */}
                        <form onSubmit={handleSubmit} className="mb-3">
                            <div className="input-group app-id-input mx-auto">
                                <input
                                    type="text"
                                    className="form-control"
                                    placeholder="Enter App ID (e.g., 12345678)"
                                    value={inputAppId}
                                    onChange={(e) => setInputAppId(e.target.value)}
                                    disabled={loading}
                                />
                                <button 
                                    className="btn btn-primary"
                                    type="submit"
                                    disabled={loading || !inputAppId.trim()}
                                >
                                    {loading ? (
                                        <>
                                            <span className="spinner-border spinner-border-sm me-2" role="status"></span>
                                            Loading...
                                        </>
                                    ) : (
                                        <>
                                            <i className="bi bi-search me-2"></i>
                                            Load Reviews
                                        </>
                                    )}
                                </button>
                            </div>
                        </form>
                    </div>

                    {error && (
                        <div className="alert alert-danger text-center" role="alert">
                            <i className="bi bi-exclamation-triangle-fill me-2"></i>
                            {error}
                        </div>
                    )}

                    {loading && (
                        <div className="d-flex justify-content-center align-items-center" style={{minHeight: '200px'}}>
                            <div className="text-center">
                                <div className="spinner-border text-primary loading-spinner" role="status">
                                    <span className="visually-hidden">Loading...</span>
                                </div>
                                <p className="mt-3 text-muted">Loading reviews for {inputAppId}...</p>
                            </div>
                        </div>
                    )}

                    {/* Reviews Content */}
                    {!loading && appId && (
                        <>
                            {/* Stats Summary */}
                            {reviews.length > 0 && (
                                <div className="card mb-4">
                                    <div className="card-body">
                                        <div className="row text-center">
                                            <div className="col">
                                                <h4 className="text-primary mb-0">{reviews.length}</h4>
                                                <small className="text-muted">Total Reviews</small>
                                            </div>
                                            <div className="col">
                                                <h4 className="text-warning mb-0">
                                                    <i className="bi bi-star-fill"></i> {(reviews.reduce((sum, r) => sum + r.rating, 0) / reviews.length).toFixed(1)}
                                                </h4>
                                                <small className="text-muted">Average Rating</small>
                                            </div>
                                            <div className="col">
                                                <div className="d-flex justify-content-center gap-1 flex-wrap">
                                                    {[5, 4, 3, 2, 1].map(stars => {
                                                        const count = reviews.filter(r => r.rating === stars).length;
                                                        return count > 0 && (
                                                            <small key={stars} className="text-muted">
                                                                {stars}★: {count}
                                                            </small>
                                                        );
                                                    })}
                                                </div>
                                                <small className="text-muted">Rating Distribution</small>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            )}

                            {/* Reviews List */}
                            <div>
                                {reviews.length === 0 ? (
                                    <div className="text-center py-5 empty-state">
                                        <i className="bi bi-chat-square-text display-4 text-muted"></i>
                                        <p className="mt-3 text-muted">No reviews found for this app.</p>
                                        <small className="text-muted">Try a different App ID</small>
                                    </div>
                                ) : (
                                    reviews.map(review => (
                                        <ReviewCard key={review.id} review={review} />
                                    ))
                                )}
                            </div>

                            {/* Refresh Button */}
                            {reviews.length > 0 && (
                                <div className="text-center mt-4">
                                    <button 
                                        className="btn btn-outline-primary"
                                        onClick={() => fetchReviews(appId)}
                                        disabled={loading}
                                    >
                                        <i className="bi bi-arrow-clockwise me-2"></i>
                                        Refresh Reviews
                                    </button>
                                </div>
                            )}
                        </>
                    )}

                    {/* Initial State */}
                    {!loading && !appId && !error && (
                        <div className="text-center py-5 empty-state">
                            <i className="bi bi-search display-4 text-muted"></i>
                            <p className="mt-3 text-muted">Enter an App ID to load reviews</p>
                            <small className="text-muted">
                                Example: com.example.app, com.company.awesomeapp
                            </small>
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
}

export default App;
