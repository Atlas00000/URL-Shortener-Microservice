package api

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

// responseWriter is a custom response writer that captures the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code before writing it
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// RateLimiter implements a simple rate limiting mechanism
type RateLimiter struct {
	requests map[string][]time.Time
	mu       sync.Mutex
	limit    int
	window   time.Duration
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

// rateLimitMiddleware implements rate limiting
func rateLimitMiddleware(limiter *RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr
			now := time.Now()

			limiter.mu.Lock()
			defer limiter.mu.Unlock()

			// Clean up old requests
			if requests, exists := limiter.requests[ip]; exists {
				var validRequests []time.Time
				for _, reqTime := range requests {
					if now.Sub(reqTime) < limiter.window {
						validRequests = append(validRequests, reqTime)
					}
				}
				limiter.requests[ip] = validRequests

				// Check if rate limit exceeded
				if len(validRequests) >= limiter.limit {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusTooManyRequests)
					json.NewEncoder(w).Encode(map[string]string{
						"error": "Rate limit exceeded",
					})
					return
				}
			}

			// Add new request
			limiter.requests[ip] = append(limiter.requests[ip], now)

			next.ServeHTTP(w, r)
		})
	}
}

// loggingMiddleware logs information about each request
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create custom response writer
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Process request
		next.ServeHTTP(rw, r)

		// Log request details
		duration := time.Since(start)
		log.Printf(
			"method=%s path=%s status=%d duration=%s ip=%s user_agent=%s",
			r.Method,
			r.URL.Path,
			rw.statusCode,
			duration,
			r.RemoteAddr,
			r.UserAgent(),
		)
	})
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// recoveryMiddleware recovers from panics and returns a 500 error
func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %v", err)

				// Set headers
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)

				// Write error response
				json.NewEncoder(w).Encode(ErrorResponse{
					Error:   "Internal Server Error",
					Message: "An unexpected error occurred",
				})
			}
		}()

		next.ServeHTTP(w, r)
	})
} 