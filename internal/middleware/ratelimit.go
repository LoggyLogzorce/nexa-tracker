package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// IPRateLimiter manages rate limiters for different IP addresses
type IPRateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	rate     rate.Limit
	burst    int
}

// NewIPRateLimiter creates a new IP-based rate limiter
func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	return &IPRateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rate:     r,
		burst:    b,
	}
}

// GetLimiter returns the rate limiter for the given IP address
func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter, exists := i.limiters[ip]
	if !exists {
		limiter = rate.NewLimiter(i.rate, i.burst)
		i.limiters[ip] = limiter
	}

	return limiter
}

// CleanupStale removes inactive limiters to free memory
func (i *IPRateLimiter) CleanupStale() {
	i.mu.Lock()
	defer i.mu.Unlock()

	for ip, limiter := range i.limiters {
		// If limiter is fully replenished (not used), remove it
		if limiter.Tokens() == float64(i.burst) {
			delete(i.limiters, ip)
		}
	}
}

// StartCleanup starts a background goroutine to periodically clean up stale limiters
func (i *IPRateLimiter) StartCleanup(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			i.CleanupStale()
		}
	}()
}

// RateLimitMiddleware creates a Gin middleware for rate limiting
func RateLimitMiddleware(limiter *IPRateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get client IP
		ip := c.ClientIP()

		// Get limiter for this IP
		ipLimiter := limiter.GetLimiter(ip)

		// Check if request is allowed
		if !ipLimiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"error":   "Too many requests. Please try again later.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// CreateRateLimiter creates a rate limiter with the specified requests per minute
func CreateRateLimiter(requestsPerMinute int) *IPRateLimiter {
	// Convert requests per minute to requests per second
	r := rate.Limit(float64(requestsPerMinute) / 60.0)

	limiter := NewIPRateLimiter(r, requestsPerMinute)

	// Start cleanup every 10 minutes
	limiter.StartCleanup(10 * time.Minute)

	return limiter
}
