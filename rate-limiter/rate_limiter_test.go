package main

import "testing"

func BenchmarkAtomicAllowParallel(b *testing.B) {
	limiter := NewRateLimiter(1_000_000_000)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			limiter.Allow()
		}
	})
}

func BenchmarkMutexAllowParallel(b *testing.B) {
	limiter := NewRateLimiterWithMutex(1_000_000_000)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			limiter.AllowWithMutex()
		}
	})
}
