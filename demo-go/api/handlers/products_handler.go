package handlers

import (
	"demo-go/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"gorm.io/gorm"
)

var (
    // Counter for order creation requests
    getAllProductsCounter = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "get_all_products_requests_total",
            Help: "Total number of get all products requests",
        },
        []string{"status"}, 
    )
)

func init() {
    // Register the counters with Prometheus
    prometheus.MustRegister(getAllProductsCounter)
}

func RegisterProductRoutes(r *gin.Engine, db *gorm.DB) {
    r.GET("/api/products", func(c *gin.Context) {
        GetAllProducts(c, db)
    })
}

func GetAllProducts(c *gin.Context, db *gorm.DB) {
    var products []models.Product

    if err := db.Find(&products).Error; err != nil {
		// Increment failure count for order retrieval
		getAllProductsCounter.WithLabelValues("failure").Inc()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching products"})
        return
    }

    // Increment success count for order retrieval
    orderRetrievalCounter.WithLabelValues("success").Inc()

    c.JSON(http.StatusOK, products)
}
