package handlers

import (
	models "demo-go/internal/models"
	services "demo-go/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

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
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // check stock
    hasStock, err := services.CheckStock(db, order.ProductID, order.Quantity)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking stock"})
        return
    }

    if !hasStock {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient stock"})
        return
    }

    var product models.Product
    if err := db.Find(&product, order.ProductID).Error; err != nil {
        c.JSON(404, gin.H{"error": "Product does not exist."})
        return
    }

    order.TotalPrice = product.Price * float64(order.Quantity)

    order.Status = "Pending"

    result := db.Create(&order)
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }

    c.JSON(http.StatusCreated, order)
}

func GetOrder(c *gin.Context, db *gorm.DB) {
    var order models.Order
    id := c.Param("id")

    if err := db.First(&order, "id = ?", id).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching order"})
        }
        return
    }

    c.JSON(http.StatusOK, order)
}