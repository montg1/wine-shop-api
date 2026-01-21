package service

import (
	"time"

	"wine-shop-api/internal/domain"
	"wine-shop-api/pkg/config"
)

type AnalyticsService struct{}

// DashboardStats contains overview statistics
type DashboardStats struct {
	TotalRevenue   float64 `json:"total_revenue"`
	TotalOrders    int64   `json:"total_orders"`
	TotalProducts  int64   `json:"total_products"`
	TotalCustomers int64   `json:"total_customers"`
}

// SalesByCategory represents sales grouped by wine category
type SalesByCategory struct {
	Category string  `json:"category"`
	Revenue  float64 `json:"revenue"`
	Count    int64   `json:"count"`
}

// TopProduct represents a top-selling product
type TopProduct struct {
	ID       uint    `json:"id"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Revenue  float64 `json:"revenue"`
}

// SalesByDay represents daily sales data
type SalesByDay struct {
	Date    string  `json:"date"`
	Revenue float64 `json:"revenue"`
	Orders  int64   `json:"orders"`
}

// RecentOrder represents a recent order for the dashboard
type RecentOrder struct {
	ID        uint      `json:"id"`
	UserEmail string    `json:"user_email"`
	Total     float64   `json:"total"`
	ItemCount int       `json:"item_count"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// GetDashboardStats returns overview statistics
func (s *AnalyticsService) GetDashboardStats() (*DashboardStats, error) {
	var stats DashboardStats

	// Total revenue from orders
	config.DB.Model(&domain.Order{}).
		Select("COALESCE(SUM(total), 0)").
		Scan(&stats.TotalRevenue)

	// Total orders
	config.DB.Model(&domain.Order{}).Count(&stats.TotalOrders)

	// Total products
	config.DB.Model(&domain.Product{}).Count(&stats.TotalProducts)

	// Total customers
	config.DB.Model(&domain.User{}).Where("role = ?", "customer").Count(&stats.TotalCustomers)

	return &stats, nil
}

// GetSalesByCategory returns sales grouped by wine category
func (s *AnalyticsService) GetSalesByCategory() ([]SalesByCategory, error) {
	var results []SalesByCategory

	config.DB.Table("order_items").
		Select("products.category, SUM(order_items.price * order_items.quantity) as revenue, COUNT(*) as count").
		Joins("JOIN products ON products.id = order_items.product_id").
		Group("products.category").
		Order("revenue DESC").
		Scan(&results)

	return results, nil
}

// GetTopProducts returns top selling products
func (s *AnalyticsService) GetTopProducts(limit int) ([]TopProduct, error) {
	var results []TopProduct

	config.DB.Table("order_items").
		Select("products.id, products.name, SUM(order_items.quantity) as quantity, SUM(order_items.price * order_items.quantity) as revenue").
		Joins("JOIN products ON products.id = order_items.product_id").
		Group("products.id, products.name").
		Order("quantity DESC").
		Limit(limit).
		Scan(&results)

	return results, nil
}

// GetSalesByDay returns daily sales for the last N days
func (s *AnalyticsService) GetSalesByDay(days int) ([]SalesByDay, error) {
	var results []SalesByDay

	startDate := time.Now().AddDate(0, 0, -days)

	config.DB.Table("orders").
		Select("DATE(created_at) as date, COALESCE(SUM(total_amount), 0) as revenue, COUNT(*) as orders").
		Where("created_at >= ?", startDate).
		Group("DATE(created_at)").
		Order("date ASC").
		Scan(&results)

	return results, nil
}

// GetRecentOrders returns the most recent orders
func (s *AnalyticsService) GetRecentOrders(limit int) ([]RecentOrder, error) {
	var orders []domain.Order
	var results []RecentOrder

	config.DB.Preload("Items").
		Order("created_at DESC").
		Limit(limit).
		Find(&orders)

	for _, order := range orders {
		// Get user email
		var user domain.User
		email := ""
		if err := config.DB.First(&user, order.UserID).Error; err == nil {
			email = user.Email
		}

		results = append(results, RecentOrder{
			ID:        order.ID,
			UserEmail: email,
			Total:     order.Total,
			ItemCount: len(order.Items),
			Status:    order.Status,
			CreatedAt: order.CreatedAt,
		})
	}

	return results, nil
}
