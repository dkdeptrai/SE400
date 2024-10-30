package main

import (
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	maxRequestsPerSecond = 2000 // Maximum number of requests per second
	server1              = "http://localhost:5500"
	server2              = "http://localhost:8090"
)

func main() {
	var wg sync.WaitGroup
	requestsPerSecond := 1 // Start with 1 request per second

	ticker := time.NewTicker(time.Second) // Create a ticker that ticks every second
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			for i := 0; i < requestsPerSecond; i++ {
				wg.Add(1)
				go func(requestNumber int) {
					defer wg.Done()
					// Alternate between servers
					var url string
					if requestNumber%2 == 0 {
						url = server1
					} else {
						url = server2
					}
					resp, err := http.Get(url + "/ping")
					if err != nil {
						log.Printf("Error sending request to %s: %v", url, err)
						return
					}
					defer resp.Body.Close()
					log.Printf("Response from %s: %s", url, resp.Status)
				}(i)
			}
			// wg.Wait() // Wait for all requests to finish

			// Increase requests per second, but cap at maxRequestsPerSecond
			if requestsPerSecond < maxRequestsPerSecond {
				requestsPerSecond += 100
			}
		}
	}
}
