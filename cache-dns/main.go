package main

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"time"
)

type Entry struct {
	IP        string
	ExpiresAt time.Time
}

func main() {
	// initialize DNS Cache
	dnsCache := NewDNSCache()

	hostNames := map[string]string{
		"google.com":        "142.250.183.14",
		"github.com":        "140.82.121.3",
		"stripe.com":        "54.187.201.55",
		"amazon.com":        "54.239.28.85",
		"netflix.com":       "52.89.124.203",
		"openai.com":        "104.18.33.45",
		"reddit.com":        "151.101.65.140",
		"stackoverflow.com": "198.252.206.1",
		"cloudflare.com":    "104.16.132.229",
		"docker.com":        "3.33.130.190",
	}

	hosts := make([]string, 0, len(hostNames))
	for host := range hostNames {
		hosts = append(hosts, host)
	}

	// seeding hostnames to DNS Cache
	for name := range hostNames {
		dnsCache.Set(name, hostNames[name], time.Duration(rand.IntN(61))*time.Second)
	}

	// 10 - reader go-routines
	var wg sync.WaitGroup
	for i := range 10 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for range 5 {
				host := hosts[rand.IntN(len(hosts))]
				ip, ok := dnsCache.Get(host)
				if ok {
					fmt.Printf("[Reader %d] HIT  %s -> %s\n", i, host, ip)
				} else {
					fmt.Printf("[Reader %d] MISS %s\n", i, host)
				}
				time.Sleep(time.Duration(rand.IntN(11)) * time.Second)
			}
		}(i)
	}

	// 1 - writer go-routine
	wg.Add(1)
	go func() {
		defer wg.Done()

		for range 5 {
			time.Sleep(10 * time.Second)

			cleaned := dnsCache.Cleanup()
			for _, host := range cleaned {
				dnsCache.Set(host, hostNames[host], time.Duration(rand.IntN(61))*time.Second)
				fmt.Printf("[Writer] refreshed %s\n", host)
			}
		}
	}()

	wg.Wait()

	fmt.Println("Simulation complete")
}
