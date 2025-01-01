package main

import (
	"bytes"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// sendReq sends a request to the server
func sendReq(m *metrics, client *http.Client, url string) {
	// Sleep to avoid sending requests at the same time.
	rn := rand.Intn(*scaleInterval)
	time.Sleep(time.Duration(rn) * time.Millisecond)

	// Get timestamp for histogram
	now := time.Now()

	// Prepare the request
	var req *http.Request
	var err error
	if *httpMethod == "POST" {
		req, err = http.NewRequest("POST", url, bytes.NewBuffer([]byte(*jsonBody)))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest("GET", url, nil)
	}

	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return
	}

	// Send the request
	res, err := client.Do(req)
	if err != nil {
		m.duration.With(prometheus.Labels{"path": url, "status": "500"}).Observe(time.Since(now).Seconds())
		log.Printf("client.Do failed: %v", err)
		return
	}
	// Read until the response is complete to reuse connection
	io.ReadAll(res.Body)

	// Close the body to reuse connection
	res.Body.Close()

	// Record request duration
	m.duration.With(prometheus.Labels{"path": url, "status": strconv.Itoa(res.StatusCode)}).Observe(time.Since(now).Seconds())
}