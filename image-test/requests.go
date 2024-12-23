package main

import (
	"bytes"
	"io"
	"log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// sendReq sends a POST request to the server with an image in the form data
func sendReq(m *metrics, client *http.Client, url string) {
	// Sleep to avoid sending requests at the same time.
	rn := rand.Intn(*scaleInterval)
	time.Sleep(time.Duration(rn) * time.Millisecond)

	// Get timestamp for request duration
	now := time.Now()

	// Open the image file
	file, err := os.Open(*imagePath)
	if err != nil {
		log.Printf("Failed to open image file: %v", err)
		m.duration.With(prometheus.Labels{"path": url, "status": "500"}).Observe(time.Since(now).Seconds())
		return
	}
	defer file.Close()

	// Prepare the body of the POST request
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", *imagePath)
	if err != nil {
		log.Printf("Failed to create form file: %v", err)
		m.duration.With(prometheus.Labels{"path": url, "status": "500"}).Observe(time.Since(now).Seconds())
		return
	}

	// Copy the image into the form file
	_, err = io.Copy(part, file)
	if err != nil {
		log.Printf("Failed to copy image data: %v", err)
		m.duration.With(prometheus.Labels{"path": url, "status": "500"}).Observe(time.Since(now).Seconds())
		return
	}

	// Close the writer to finalize the body
	err = writer.Close()
	if err != nil {
		log.Printf("Failed to close multipart writer: %v", err)
		m.duration.With(prometheus.Labels{"path": url, "status": "500"}).Observe(time.Since(now).Seconds())
		return
	}

	// Send the POST request
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		m.duration.With(prometheus.Labels{"path": url, "status": "500"}).Observe(time.Since(now).Seconds())
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

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
