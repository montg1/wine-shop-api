package repository

import (
	"wine-shop-api/internal/domain"

	"gorm.io/gorm"
)

type cartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) CartRepository {
	return &cartRepository{db: db}
}

func (r *cartRepository) FindByUserID(userID uint) (*domain.Cart, error) {
	var cart domain.Cart
	err := r.db.Preload("Items.Product").Where("user_id = ?", userID).First(&cart).Error
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r *cartRepository) Create(cart *domain.Cart) error {
	return r.db.Create(cart).Error
}

func (r *cartRepository) FindCartItem(cartID, productID uint) (*domain.CartItem, error) {
	var item domain.CartItem
	err := r.db.Where("cart_id = ? AND product_id = ?", cartID, productID).First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *cartRepository) SaveCartItem(item *domain.CartItem) error {
	return r.db.Save(item).Error
}

func (r *cartRepository) CreateCartItem(item *domain.CartItem) error {
	return r.db.Create(item).Error
}

func (r *cartRepository) DeleteCartItems(cartID uint) error {
	return r.db.Where("cart_id = ?", cartID).Delete(&domain.CartItem{}).Error
}
