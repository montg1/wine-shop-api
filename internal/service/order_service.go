package service

import (
	"errors"

	"wine-shop-api/internal/domain"
	"wine-shop-api/pkg/config"

	"gorm.io/gorm"
)

type OrderService struct {
	CartService *CartService
}

func (s *OrderService) CreateOrder(userID uint) (*domain.Order, error) {
	// 1. Get Cart
	cart, err := s.CartService.GetCart(userID)
	if err != nil {
		return nil, err
	}

	if len(cart.Items) == 0 {
		return nil, errors.New("cart is empty")
	}

	// 2. Calculate Total and Create Order Items
	var total float64
	var orderItems []domain.OrderItem

	for _, item := range cart.Items {
		total += item.Product.Price * float64(item.Quantity)
		orderItems = append(orderItems, domain.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Product.Price, // Snapshot price at purchase time
		})
	}

	// 3. Create Order
	order := domain.Order{
		UserID: userID,
		Total:  total,
		Status: "Paid", // Simplified for this demo
		Items:  orderItems,
	}

	tx := config.DB.Begin()

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 4. Clear Cart
	if err := tx.Where("cart_id = ?", cart.ID).Delete(&domain.CartItem{}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 5. Update Stock (Optional but good practice)
	for _, item := range cart.Items {
		if err := tx.Model(&domain.Product{}).Where("id = ?", item.ProductID).UpdateColumn("stock", gorm.Expr("stock - ?", item.Quantity)).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	tx.Commit()

	return &order, nil
}

func (s *OrderService) GetOrders(userID uint) ([]domain.Order, error) {
	var orders []domain.Order
	if err := config.DB.Preload("Items.Product").Where("user_id = ?", userID).Order("created_at desc").Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}
