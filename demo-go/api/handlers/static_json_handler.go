package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var (
    // Counter for get static json request
    getStaticJsonCounter = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "get_static_json_requests_total",
            Help: "Total number of get static json requests",
        },
        []string{"status"}, 
    )
)

func init() {
    // Register the counter with Prometheus
    prometheus.MustRegister(getStaticJsonCounter)
}

func RegisterStaticJsonRoutes(r *gin.Engine) {
    r.GET("/api/static-json", func(c *gin.Context) {
        GetStaticJson(c)
    })
}

func GetStaticJson(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello, World!",
		"status":  "success",
		"data": gin.H{
			"user": gin.H{
				"id":    1,
				"name":  "John Doe",
				"email": "john.doe@example.com",
			},
			"posts": []gin.H{
				{
					"id":      101,
					"title":   "First Post",
					"content": "This is the content of the first post.",
				},
				{
					"id":      102,
					"title":   "Second Post",
					"content": "This is the content of the second post.",
				},
			},
		},
	})

    // Increment success count for getting static json
    getStaticJsonCounter.WithLabelValues("success").Inc()
}
