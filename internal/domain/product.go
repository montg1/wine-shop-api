package domain

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	ImageURL    string  `json:"image_url"`
	Category    string  `json:"category"` // e.g., "Red", "White", "Sparkling"
}

// IsValid validates the product fields
func (p *Product) IsValid() bool {
	if p.Name == "" {
		return false
	}
	if p.Price < 0 {
		return false
	}
	if p.Stock < 0 {
		return false
	}
	return true
}

// IsValidCategory checks if the category is valid
func (p *Product) IsValidCategory() bool {
	validCategories := []string{"Red", "White", "Rosé", "Sparkling", "Dessert"}
	for _, cat := range validCategories {
		if p.Category == cat {
			return true
		}
	}
	return false
}
