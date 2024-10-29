package models

type Order struct {
    ID         uint    `json:"id" gorm:"primaryKey"`
    ProductID  uint    `json:"product_id"`
    Product    Product `json:"product" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
    Quantity   int     `json:"quantity" gorm:"not null;check:quantity > 0"`
    TotalPrice float64 `json:"total_price" gorm:"not null;check:total_price >= 0"`
    Status     string  `json:"status"`
}