package repository

import (
	"time"

	"wine-shop-api/internal/domain"

	"gorm.io/gorm"
)

type analyticsRepository struct {
	db *gorm.DB
}

func NewAnalyticsRepository(db *gorm.DB) AnalyticsRepository {
	return &analyticsRepository{db: db}
}

func (r *analyticsRepository) GetTotalRevenue() (float64, error) {
	var revenue float64
	err := r.db.Model(&domain.Order{}).
		Select("COALESCE(SUM(total), 0)").
		Scan(&revenue).Error
	return revenue, err
}

func (r *analyticsRepository) GetTotalOrders() (int64, error) {
	var count int64
	err := r.db.Model(&domain.Order{}).Count(&count).Error
	return count, err
}

func (r *analyticsRepository) GetTotalProducts() (int64, error) {
	var count int64
	err := r.db.Model(&domain.Product{}).Count(&count).Error
	return count, err
}

func (r *analyticsRepository) GetTotalCustomers() (int64, error) {
	var count int64
	err := r.db.Model(&domain.User{}).Where("role = ?", "customer").Count(&count).Error
	return count, err
}

func (r *analyticsRepository) GetSalesByCategory() ([]SalesByCategory, error) {
	var results []SalesByCategory
	err := r.db.Table("order_items").
		Select("products.category, SUM(order_items.price * order_items.quantity) as revenue, COUNT(*) as count").
		Joins("JOIN products ON products.id = order_items.product_id").
		Group("products.category").
		Order("revenue DESC").
		Scan(&results).Error
	return results, err
}

func (r *analyticsRepository) GetTopProducts(limit int) ([]TopProduct, error) {
	var results []TopProduct
	err := r.db.Table("order_items").
		Select("products.id, products.name, SUM(order_items.quantity) as quantity, SUM(order_items.price * order_items.quantity) as revenue").
		Joins("JOIN products ON products.id = order_items.product_id").
		Group("products.id, products.name").
		Order("quantity DESC").
		Limit(limit).
		Scan(&results).Error
	return results, err
}

func (r *analyticsRepository) GetSalesByDay(days int) ([]SalesByDay, error) {
	results := []SalesByDay{}
	startDate := time.Now().AddDate(0, 0, -days)

	err := r.db.Table("orders").
		Select("TO_CHAR(created_at, 'YYYY-MM-DD') as date, COALESCE(SUM(total), 0) as revenue, COUNT(*) as orders").
		Where("created_at >= ?", startDate).
		Group("TO_CHAR(created_at, 'YYYY-MM-DD')").
		Order("date ASC").
		Scan(&results).Error
	if err != nil {
		return []SalesByDay{}, err
	}
	return results, nil
}

func (r *analyticsRepository) GetRecentOrders(limit int) ([]RecentOrderRow, error) {
	var orders []domain.Order
	r.db.Preload("Items").
		Order("created_at DESC").
		Limit(limit).
		Find(&orders)

	var results []RecentOrderRow
	for _, order := range orders {
		var user domain.User
		email := ""
		if err := r.db.First(&user, order.UserID).Error; err == nil {
			email = user.Email
		}

		results = append(results, RecentOrderRow{
			ID:        order.ID,
			UserID:    order.UserID,
			Email:     email,
			Total:     order.Total,
			ItemCount: len(order.Items),
			Status:    order.Status,
			CreatedAt: order.CreatedAt,
		})
	}

	return results, nil
}
