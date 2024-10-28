package services

import (
	"gorm.io/gorm"

	models "demo-go/internal/models"
)

func CheckStock(db *gorm.DB, productID uint, quantity int) (bool, error) {
    var product models.Product
    if err := db.First(&product, productID).Error; err != nil {
        return false, err // Product does not exists
    }

    if product.Stock < quantity {
        return false, nil
    }

    // Get from stock
    product.Stock -= quantity
    if err := db.Save(&product).Error; err != nil {
        return false, err
    }

    return true, nil 
}
