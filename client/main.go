package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Define flag variables
var (
	maxClients    = flag.Int("maxClients", 100, "Maximum number of virtual clients")
	scaleInterval = flag.Int("scaleInterval", 1, "Scale interval in milliseconds")
	randomSleep   = flag.Int("randomSleep", 1000, "Random sleep from 0 to target microseconds")
	port          = flag.Int("port", 8081, "Port for Prometheus HTTP server metrics")
	targets       multiString // Custom flag type to handle multiple targets
)

func init() {
	// Register the custom flag for multiple targets
	flag.Var(&targets, "target", "Target URL for the servers; can be specified multiple times for different URLs")
}

// multiString is a custom flag type to handle multiple string values for -target
type multiString []string

func (ms *multiString) String() string {
	return "multiple target URLs"
}

func (ms *multiString) Set(value string) error {
	*ms = append(*ms, value)
	return nil
}

func main() {
	// Sleep for 5 seconds before running test
	time.Sleep(5 * time.Second)

	// Parse the command line into the defined flags
	flag.Parse()

	// Create Prometheus registry
	reg := prometheus.NewRegistry()
	m := NewMetrics(reg)

	// Create Prometheus HTTP server to expose metrics
	pMux := http.NewServeMux()
	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
	pMux.Handle("/metrics", promHandler)

	go func() {
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), pMux))
	}()

	// Create transport and client to reuse connection pool
	tr := &http.Transport{}
	client := &http.Client{Transport: tr}

	// Create job queue
	var ch = make(chan string, *maxClients*2)
	var wg sync.WaitGroup

	// Slowly increase the number of virtual clients
	for clients := 0; clients <= *maxClients; clients++ {
		wg.Add(1)

		for i := 0; i < clients; i++ {
			go func() {
				for {
					url, ok := <-ch
					if !ok {
						wg.Done()
						return
					}
					sendReq(m, client, url)
				}
			}()
		}

		doWork(ch, clients)

		time.Sleep(time.Duration(*scaleInterval) * time.Millisecond)
	}
}

func doWork(ch chan string, clients int) {
	if clients == *maxClients {
		for {
			url := getRandomTarget()
			ch <- url
			sleep(*randomSleep)
		}
	}

	for i := 0; i < clients; i++ {
		url := getRandomTarget()
		ch <- url
		sleep(*randomSleep)
	}
}

func getRandomTarget() string {
	if len(targets) == 0 {
		return ""
	}
	return targets[rand.Intn(len(targets))]
}

func sleep(us int) {
	r := rand.Intn(us)
	time.Sleep(time.Duration(r) * time.Microsecond)
}
