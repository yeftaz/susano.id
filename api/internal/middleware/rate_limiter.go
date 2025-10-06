package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/yeftaz/susano.id/api/internal/config"
	"github.com/yeftaz/susano.id/api/pkg/response"
)

type visitor struct {
	limiter  *rateLimiter
	lastSeen time.Time
}

type rateLimiter struct {
	tokens     int
	maxTokens  int
	refillRate time.Duration
	lastRefill time.Time
	mu         sync.Mutex
}

func newRateLimiter(maxTokens int, refillRate time.Duration) *rateLimiter {
	return &rateLimiter{
		tokens:     maxTokens,
		maxTokens:  maxTokens,
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

func (rl *rateLimiter) allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(rl.lastRefill)

	// Refill tokens based on time elapsed
	if elapsed >= rl.refillRate {
		rl.tokens = rl.maxTokens
		rl.lastRefill = now
	}

	// Check if request is allowed
	if rl.tokens > 0 {
		rl.tokens--
		return true
	}

	return false
}

var (
	visitors = make(map[string]*visitor)
	mu       sync.Mutex
)

// RateLimiter middleware limits requests per IP address
func RateLimiter(cfg *config.Config) func(http.Handler) http.Handler {
	// Cleanup old visitors every minute
	go cleanupVisitors()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := getIP(r)

			mu.Lock()
			v, exists := visitors[ip]
			if !exists {
				v = &visitor{
					limiter: newRateLimiter(cfg.RateLimitRequests, cfg.RateLimitWindow),
				}
				visitors[ip] = v
			}
			v.lastSeen = time.Now()
			mu.Unlock()

			if !v.limiter.allow() {
				response.Error(w, http.StatusTooManyRequests, "Rate limit exceeded")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func getIP(r *http.Request) string {
	// Check X-Forwarded-For header
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		// Take the first IP if multiple are present
		ips := strings.Split(forwarded, ",")
		return strings.TrimSpace(ips[0])
	}

	// Check X-Real-IP header
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Fall back to RemoteAddr
	return r.RemoteAddr
}

func cleanupVisitors() {
	for {
		time.Sleep(1 * time.Minute)

		mu.Lock()
		for ip, v := range visitors {
			if time.Since(v.lastSeen) > 3*time.Minute {
				delete(visitors, ip)
			}
		}
		mu.Unlock()
	}
}
