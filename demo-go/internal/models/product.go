package models

type Product struct {
    ID    uint    `json:"id" gorm:"primaryKey"`
    Name  string  `json:"name" gorm:"not null"`
    Price float64 `json:"price" gorm:"not null; check: price > 0"`
    Stock int     `json:"stock" gorm:"not null; check: stock > 0"` 
}
