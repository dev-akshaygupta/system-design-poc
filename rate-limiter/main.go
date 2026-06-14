package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// with atomic
func main() {
	rateLimiter := NewRateLimiter(10)

	var wg sync.WaitGroup
	done := make(chan struct{})
	wg.Go(func() {
		rateLimiter.StartResetTicker(done)
	})

	for batch := range 2 {
		var allowed int64
		var denied int64

		var bwg sync.WaitGroup
		for range 25 {
			bwg.Go(func() {
				if rateLimiter.Allow() {
					atomic.AddInt64(&allowed, 1)
				} else {
					atomic.AddInt64(&denied, 1)
				}
			})
		}
		bwg.Wait()
		fmt.Printf("Batch %d -> Allowed: %d, Denied: %d\n", batch, allowed, denied)
		time.Sleep(1200 * time.Millisecond)
	}
	close(done)
	wg.Wait()
}

// with mutex
// func main() {
// 	rateLimiterWithMU := NewRateLimiterWithMutex(10)

// 	var wg sync.WaitGroup
// 	done := make(chan struct{})
// 	wg.Go(func() {
// 		rateLimiterWithMU.StartResetTickerWithMutex(done)
// 	})

// 	for batch := range 2 {
// 		var allowed int64
// 		var denied int64

// 		var bwg sync.WaitGroup
// 		for range 25 {
// 			bwg.Go(func() {
// 				if rateLimiterWithMU.AllowWithMutex() {
// 					atomic.AddInt64(&allowed, 1)
// 				} else {
// 					atomic.AddInt64(&denied, 1)
// 				}
// 			})
// 		}
// 		bwg.Wait()
// 		fmt.Printf("Batch %d -> Allowed: %d, Denied: %d\n", batch, allowed, denied)
// 		time.Sleep(1200 * time.Millisecond)
// 	}
// 	close(done)
// 	wg.Wait()
// }
