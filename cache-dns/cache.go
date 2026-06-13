package main

import (
	"fmt"
	"sync"
	"time"
)

type DNSCache struct {
	mu    sync.RWMutex
	store map[string]Entry
}

func NewDNSCache() *DNSCache {
	dnsCache := DNSCache{
		store: make(map[string]Entry),
	}
	return &dnsCache
}

func (dnsc *DNSCache) Set(host string, ip string, ttl time.Duration) {
	dnsc.mu.Lock()
	defer dnsc.mu.Unlock()

	dnsc.store[host] = Entry{IP: ip, ExpiresAt: time.Now().Add(ttl)}
}

func (dnsc *DNSCache) Get(host string) (string, bool) {
	dnsc.mu.RLock()
	defer dnsc.mu.RUnlock()

	entry, ok := dnsc.store[host]
	if !ok {
		return "", false
	}

	if time.Now().After(entry.ExpiresAt) {
		return "", false
	}
	return entry.IP, true
}

func (dnsc *DNSCache) Cleanup() []string {
	dnsc.mu.Lock()
	defer dnsc.mu.Unlock()

	cleaned := []string{}
	count := 0
	for key, entry := range dnsc.store {
		if time.Now().After(entry.ExpiresAt) {
			delete(dnsc.store, key)
			count++
			cleaned = append(cleaned, key)
		}
	}

	fmt.Println("Cleanup Done. Total keys cleaned:", count)
	return cleaned
}
