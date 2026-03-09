package repository

import (
	"wine-shop-api/internal/domain"

	"gorm.io/gorm"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	FindByEmail(email string) (*domain.User, error)
	Create(user *domain.User) error
	FindByID(id uint) (*domain.User, error)
	Save(user *domain.User) error
}

// ProductRepository defines the interface for product data access
type ProductRepository interface {
	Create(product *domain.Product) error
	FindAll(offset, limit int, search, category string) ([]domain.Product, int64, error)
	FindByID(id uint) (*domain.Product, error)
	Save(product *domain.Product) error
	Delete(id uint) error
}

// CartRepository defines the interface for cart data access
type CartRepository interface {
	FindByUserID(userID uint) (*domain.Cart, error)
	Create(cart *domain.Cart) error
	FindCartItem(cartID, productID uint) (*domain.CartItem, error)
	SaveCartItem(item *domain.CartItem) error
	CreateCartItem(item *domain.CartItem) error
	DeleteCartItems(cartID uint) error
}

// OrderRepository defines the interface for order data access
type OrderRepository interface {
	CreateOrder(order *domain.Order) error
	DeleteCartItemsInTx(tx *gorm.DB, cartID uint) error
	UpdateStockInTx(tx *gorm.DB, productID uint, quantity int) error
	FindByUserID(userID uint) ([]domain.Order, error)
	BeginTx() *gorm.DB
}

// ReviewRepository defines the interface for review data access
type ReviewRepository interface {
	FindByProductAndUser(productID, userID uint) (*domain.Review, error)
	Create(review *domain.Review) error
	PreloadUser(review *domain.Review) error
	FindByProductID(productID uint) ([]domain.Review, error)
	AverageRating(productID uint) (float64, int64, error)
	FindByID(id uint) (*domain.Review, error)
	Delete(review *domain.Review) error
}

// AnalyticsRepository defines the interface for analytics data access
type AnalyticsRepository interface {
	GetTotalRevenue() (float64, error)
	GetTotalOrders() (int64, error)
	GetTotalProducts() (int64, error)
	GetTotalCustomers() (int64, error)
	GetSalesByCategory() ([]SalesByCategory, error)
	GetTopProducts(limit int) ([]TopProduct, error)
	GetSalesByDay(days int) ([]SalesByDay, error)
	GetRecentOrders(limit int) ([]RecentOrderRow, error)
}

// Analytics DTOs used by repository layer
type SalesByCategory struct {
	Category string  `json:"category"`
	Revenue  float64 `json:"revenue"`
	Count    int64   `json:"count"`
}

type TopProduct struct {
	ID       uint    `json:"id"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Revenue  float64 `json:"revenue"`
}

type SalesByDay struct {
	Date    string  `json:"date"`
	Revenue float64 `json:"revenue"`
	Orders  int64   `json:"orders"`
}

type RecentOrderRow struct {
	ID        uint
	UserID    uint
	Email     string
	Total     float64
	ItemCount int
	Status    string
	CreatedAt interface{}
}
