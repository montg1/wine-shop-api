package repository

import (
	"wine-shop-api/internal/domain"

	"gorm.io/gorm"
)

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) BeginTx() *gorm.DB {
	return r.db.Begin()
}

func (r *orderRepository) CreateOrder(order *domain.Order) error {
	return r.db.Create(order).Error
}

func (r *orderRepository) DeleteCartItemsInTx(tx *gorm.DB, cartID uint) error {
	return tx.Where("cart_id = ?", cartID).Delete(&domain.CartItem{}).Error
}

func (r *orderRepository) UpdateStockInTx(tx *gorm.DB, productID uint, quantity int) error {
	return tx.Model(&domain.Product{}).Where("id = ?", productID).
		UpdateColumn("stock", gorm.Expr("stock - ?", quantity)).Error
}

func (r *orderRepository) FindByUserID(userID uint) ([]domain.Order, error) {
	var orders []domain.Order
	if err := r.db.Preload("Items.Product").Where("user_id = ?", userID).
		Order("created_at desc").Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}
