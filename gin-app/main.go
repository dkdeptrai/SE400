package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

var (
	pingRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "go_ping_requests_total",
			Help: "Total number of requests to the /ping endpoint",
		},
		[]string{"status"},
	)

	cpuUsage = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "go_cpu_usage_percentage",
			Help: "Current CPU usage as a percentage",
		},
	)

	ramUsage = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "go_ram_usage_bytes",
			Help: "Current RAM usage in bytes",
		},
	)

	ramUsagePercentage = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "go_ram_usage_percentage",
			Help: "Current RAM usage as a percentage",
		},
	)

	requestLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "go_request_latency_seconds",
			Help:    "Histogram of latencies for requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path"},
	)
)

func recordSystemMetrics() {
	for {
		cpuPercent, _ := cpu.Percent(0, false)
		if len(cpuPercent) > 0 {
			cpuUsage.Set(cpuPercent[0])
		}

		vMem, _ := mem.VirtualMemory()
		ramUsage.Set(float64(vMem.Used))

		ramUsagePercentage.Set(vMem.UsedPercent)

		time.Sleep(5 * time.Second) // Poll every 5 seconds
	}
}

func main() {
	// Register Prometheus metrics
	prometheus.MustRegister(pingRequests, cpuUsage, ramUsage, ramUsagePercentage, requestLatency)

	// Run system metrics recording in a goroutine
	go recordSystemMetrics()

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		// Measure request latency
		start := time.Now()
		defer func() {
			duration := time.Since(start).Seconds()
			requestLatency.WithLabelValues("/ping").Observe(duration)
		}()

		pingRequests.WithLabelValues("200").Inc()
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	r.Run(":8090")
}
