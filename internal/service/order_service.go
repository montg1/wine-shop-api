package service

import (
	"errors"

	"wine-shop-api/internal/domain"
	"wine-shop-api/internal/repository"
)

type OrderService struct {
	Repo        repository.OrderRepository
	CartService *CartService
}

func NewOrderService(repo repository.OrderRepository, cartService *CartService) *OrderService {
	return &OrderService{Repo: repo, CartService: cartService}
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
			Price:     item.Product.Price,
		})
	}

	// 3. Create Order
	order := domain.Order{
		UserID: userID,
		Total:  total,
		Status: "Paid",
		Items:  orderItems,
	}

	tx := s.Repo.BeginTx()

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 4. Clear Cart
	if err := s.Repo.DeleteCartItemsInTx(tx, cart.ID); err != nil {
		tx.Rollback()
		return nil, err
	}

	// 5. Update Stock
	for _, item := range cart.Items {
		if err := s.Repo.UpdateStockInTx(tx, item.ProductID, item.Quantity); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	tx.Commit()

	return &order, nil
}

func (s *OrderService) GetOrders(userID uint) ([]domain.Order, error) {
	return s.Repo.FindByUserID(userID)
}
