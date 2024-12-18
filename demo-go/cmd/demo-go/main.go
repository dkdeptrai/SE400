package main

import (
	"demo-go/api/handlers"
	"demo-go/internal/config"
	"demo-go/internal/database"
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
    // Load configuration
    cfg := config.LoadConfig()

    // Initialize database
    db := database.InitDB(cfg.Database.DSN)

    // Register Prometheus metrics
	prometheus.MustRegister(pingRequests, cpuUsage, ramUsage, ramUsagePercentage, requestLatency)

	// Start recording system metrics in a separate goroutine
	go recordSystemMetrics()

    // Set up the Gin router
    r := gin.Default()

    // Middleware to measure request latency
	r.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start).Seconds()
		requestLatency.WithLabelValues(c.FullPath()).Observe(duration)
	})

	handlers.ConvertToMonochrome()

    // Custom ping endpoint to demonstrate metrics tracking
	// r.GET("/ping", func(c *gin.Context) {
	// 	pingRequests.WithLabelValues("200").Inc()
	// 	c.JSON(200, gin.H{
	// 		"message": "pong",
	// 	})
	// })

    // Register routes
    handlers.RegisterRoutes(r, db)

	// Prometheus metrics endpoint
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

    // Start the server
    r.Run(cfg.Server.Address)
}
