package middleware

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"pkl-attendance/backend/utils"
)

type limiterEntry struct {
	count   int
	resetAt time.Time
}

type RateLimiter struct {
	mu      sync.Mutex
	entries map[string]*limiterEntry
	limit   int
	window  time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		entries: make(map[string]*limiterEntry),
		limit:   limit,
		window:  window,
	}
	go rl.cleanupLoop()
	return rl
}

func (rl *RateLimiter) cleanupLoop() {
	t := time.NewTicker(time.Minute)
	defer t.Stop()
	for range t.C {
		now := time.Now()
		rl.mu.Lock()
		for k, v := range rl.entries {
			if now.After(v.resetAt) {
				delete(rl.entries, k)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.ClientIP()
		rl.mu.Lock()
		now := time.Now()
		e, ok := rl.entries[key]
		if !ok || now.After(e.resetAt) {
			e = &limiterEntry{count: 0, resetAt: now.Add(rl.window)}
			rl.entries[key] = e
		}
		e.count++
		over := e.count > rl.limit
		rl.mu.Unlock()
		if over {
			utils.Error(c, 429, "too many requests, please try again later")
			c.Abort()
			return
		}
		c.Next()
	}
}