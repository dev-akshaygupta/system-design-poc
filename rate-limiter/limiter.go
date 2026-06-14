package main

import (
	"sync"
	"sync/atomic"
	"time"
)

type RateLimiter struct {
	counter atomic.Int64
	limit   int64
}

type RateLimiterWithMutex struct {
	mu      sync.Mutex
	counter int64
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

func NewRateLimiterWithMutex(limit int64) *RateLimiterWithMutex {
	ratelimiter := RateLimiterWithMutex{
		limit: limit,
	}
	return &ratelimiter
}

func (rlm *RateLimiterWithMutex) AllowWithMutex() bool {
	rlm.mu.Lock()
	rlm.counter++
	if rlm.counter <= rlm.limit {
		return true
	}
	rlm.mu.Unlock()
	return false
}

func (rlm *RateLimiterWithMutex) ResetWithMutex() {
	rlm.mu.Lock()
	defer rlm.mu.Unlock()
	rlm.counter = 0
}

func (rlm *RateLimiterWithMutex) StartResetTickerWithMutex(done <-chan struct{}) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rlm.ResetWithMutex()
		case <-done:
			return
		}
	}
}
