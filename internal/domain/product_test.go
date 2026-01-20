package domain

import (
	"testing"
)

func TestProduct_Validation(t *testing.T) {
	tests := []struct {
		name    string
		product Product
		valid   bool
	}{
		{
			name: "Valid product",
			product: Product{
				Name:     "Pinot Noir",
				Price:    45.00,
				Stock:    10,
				Category: "Red",
			},
			valid: true,
		},
		{
			name: "Empty name",
			product: Product{
				Name:  "",
				Price: 45.00,
				Stock: 10,
			},
			valid: false,
		},
		{
			name: "Negative price",
			product: Product{
				Name:  "Test Wine",
				Price: -10.00,
				Stock: 10,
			},
			valid: false,
		},
		{
			name: "Negative stock",
			product: Product{
				Name:  "Test Wine",
				Price: 45.00,
				Stock: -5,
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.product.IsValid()
			if result != tt.valid {
				t.Errorf("Product.IsValid() = %v, want %v", result, tt.valid)
			}
		})
	}
}

func TestProduct_Categories(t *testing.T) {
	validCategories := []string{"Red", "White", "Rosé", "Sparkling", "Dessert"}

	for _, cat := range validCategories {
		p := Product{
			Name:     "Test Wine",
			Price:    50.00,
			Stock:    10,
			Category: cat,
		}
		if !p.IsValidCategory() {
			t.Errorf("Category %s should be valid", cat)
		}
	}

	// Test invalid category
	p := Product{
		Name:     "Test Wine",
		Price:    50.00,
		Stock:    10,
		Category: "InvalidCategory",
	}
	if p.IsValidCategory() {
		t.Error("Invalid category should return false")
	}
}
