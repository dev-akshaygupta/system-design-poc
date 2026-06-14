package main

import (
	"sync"
	"sync/atomic"
	"time"
)

type RateLimiter struct {
	mu      sync.RWMutex
	counter atomic.Int64
	limit   int64
}

func NewRateLimiter(limit int64) *RateLimiter {
	ratelimiter := RateLimiter{
		limit: limit,
	}
	return &ratelimiter
}

func (rl *RateLimiter) Allow() bool {
	count := rl.counter.Add(1)
	if count <= rl.limit {
		return true
	}
	return false
}

func (rl *RateLimiter) Reset() {
	rl.counter.Store(0)
}

func (rl *RateLimiter) StartResetTicker(done <-chan struct{}) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rl.Reset()
		case <-done:
			return
		}
	}
}
