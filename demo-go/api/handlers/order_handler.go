package handlers

import (
	"demo-go/internal/models"
	"demo-go/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"gorm.io/gorm"
)

var (
    // Counter for order creation requests
    orderCreationCounter = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "order_creation_requests_total",
            Help: "Total number of order creation requests",
        },
        []string{"status"}, // Track success/failure
    )

    // Counter for order retrieval requests
    orderRetrievalCounter = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "order_retrieval_requests_total",
            Help: "Total number of order retrieval requests",
        },
        []string{"status"}, // Track success/failure
    )
)

func init() {
    // Register the counters with Prometheus
    prometheus.MustRegister(orderCreationCounter)
    prometheus.MustRegister(orderRetrievalCounter)
}

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
    r.POST("/orders", func(c *gin.Context) {
        CreateOrder(c, db)
    })

    r.GET("/orders/:id", func(c *gin.Context) {
        GetOrder(c, db)
    })
}

func CreateOrder(c *gin.Context, db *gorm.DB) {
    var order models.Order
    if err := c.ShouldBindJSON(&order); err != nil {
        // Increment failure count for order creation
        orderCreationCounter.WithLabelValues("failure").Inc()
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Check stock
    hasStock, err := services.CheckStock(db, order.ProductID, order.Quantity)
    if err != nil {
        // Increment failure count for order creation
        orderCreationCounter.WithLabelValues("failure").Inc()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking stock"})
        return
    }

    if !hasStock {
        // Increment failure count for order creation
        orderCreationCounter.WithLabelValues("failure").Inc()
        c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient stock"})
        return
    }

    var product models.Product
    if err := db.Find(&product, order.ProductID).Error; err != nil {
        // Increment failure count for order creation
        orderCreationCounter.WithLabelValues("failure").Inc()
        c.JSON(404, gin.H{"error": "Product does not exist."})
        return
    }

    order.TotalPrice = product.Price * float64(order.Quantity)
    order.Status = "Pending"

    result := db.Create(&order)
    if result.Error != nil {
        // Increment failure count for order creation
        orderCreationCounter.WithLabelValues("failure").Inc()
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }

    // Increment success count for order creation
    orderCreationCounter.WithLabelValues("success").Inc()

    c.JSON(http.StatusCreated, order)
}

func GetOrder(c *gin.Context, db *gorm.DB) {
    var order models.Order
    id := c.Param("id")

    if err := db.First(&order, "id = ?", id).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            // Increment failure count for order retrieval
            orderRetrievalCounter.WithLabelValues("failure").Inc()
            c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
        } else {
            // Increment failure count for order retrieval
            orderRetrievalCounter.WithLabelValues("failure").Inc()
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching order"})
        }
        return
    }

    // Increment success count for order retrieval
    orderRetrievalCounter.WithLabelValues("success").Inc()

    c.JSON(http.StatusOK, order)
}
