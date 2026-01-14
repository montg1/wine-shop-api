package service

import (
	"errors"

	"wine-shop-api/internal/domain"
	"wine-shop-api/pkg/config"

	"gorm.io/gorm"
)

type CartService struct{}

func (s *CartService) GetCart(userID uint) (*domain.Cart, error) {
	var cart domain.Cart
	// Find cart for user, preload items and their products
	err := config.DB.Preload("Items.Product").Where("user_id = ?", userID).First(&cart).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create a new cart if one doesn't exist
			cart = domain.Cart{UserID: userID}
			if err := config.DB.Create(&cart).Error; err != nil {
				return nil, err
			}
			return &cart, nil
		}
		return nil, err
	}
	return &cart, nil
}

func (s *CartService) AddToCart(userID uint, productID uint, quantity int) error {
	cart, err := s.GetCart(userID)
	if err != nil {
		return err
	}

	// Check if product exists
	var product domain.Product
	if err := config.DB.First(&product, productID).Error; err != nil {
		return errors.New("product not found")
	}

	// Check if item already exists in cart
	var cartItem domain.CartItem
	err = config.DB.Where("cart_id = ? AND product_id = ?", cart.ID, productID).First(&cartItem).Error

	if err == nil {
		// Update quantity
		cartItem.Quantity += quantity
		return config.DB.Save(&cartItem).Error
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		// Create new item
		newItem := domain.CartItem{
			CartID:    cart.ID,
			ProductID: productID,
			Quantity:  quantity,
		}
		return config.DB.Create(&newItem).Error
	}

	return err
}

func (s *CartService) ClearCart(userID uint) error {
	cart, err := s.GetCart(userID)
	if err != nil {
		return err
	}

	// Delete all items in the cart
	return config.DB.Where("cart_id = ?", cart.ID).Delete(&domain.CartItem{}).Error
}
