package service

import (
	"time"

	"wine-shop-api/internal/repository"
)

type AnalyticsService struct {
	Repo repository.AnalyticsRepository
}

func NewAnalyticsService(repo repository.AnalyticsRepository) *AnalyticsService {
	return &AnalyticsService{Repo: repo}
}

// DashboardStats contains overview statistics
type DashboardStats struct {
	TotalRevenue   float64 `json:"total_revenue"`
	TotalOrders    int64   `json:"total_orders"`
	TotalProducts  int64   `json:"total_products"`
	TotalCustomers int64   `json:"total_customers"`
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

	revenue, err := s.Repo.GetTotalRevenue()
	if err != nil {
		return nil, err
	}
	stats.TotalRevenue = revenue

	orders, err := s.Repo.GetTotalOrders()
	if err != nil {
		return nil, err
	}
	stats.TotalOrders = orders

	products, err := s.Repo.GetTotalProducts()
	if err != nil {
		return nil, err
	}
	stats.TotalProducts = products

	customers, err := s.Repo.GetTotalCustomers()
	if err != nil {
		return nil, err
	}
	stats.TotalCustomers = customers

	return &stats, nil
}

// GetSalesByCategory returns sales grouped by wine category
func (s *AnalyticsService) GetSalesByCategory() ([]repository.SalesByCategory, error) {
	return s.Repo.GetSalesByCategory()
}

// GetTopProducts returns top selling products
func (s *AnalyticsService) GetTopProducts(limit int) ([]repository.TopProduct, error) {
	return s.Repo.GetTopProducts(limit)
}

// GetSalesByDay returns daily sales for the last N days
func (s *AnalyticsService) GetSalesByDay(days int) ([]repository.SalesByDay, error) {
	return s.Repo.GetSalesByDay(days)
}

// GetRecentOrders returns the most recent orders
func (s *AnalyticsService) GetRecentOrders(limit int) ([]RecentOrder, error) {
	rows, err := s.Repo.GetRecentOrders(limit)
	if err != nil {
		return nil, err
	}

	var results []RecentOrder
	for _, row := range rows {
		var createdAt time.Time
		if t, ok := row.CreatedAt.(time.Time); ok {
			createdAt = t
		}

		results = append(results, RecentOrder{
			ID:        row.ID,
			UserEmail: row.Email,
			Total:     row.Total,
			ItemCount: row.ItemCount,
			Status:    row.Status,
			CreatedAt: createdAt,
		})
	}

	return results, nil
}
