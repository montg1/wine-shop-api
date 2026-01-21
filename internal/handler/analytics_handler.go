package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"wine-shop-api/internal/service"
)

type AnalyticsHandler struct {
	Service *service.AnalyticsService
}

// GetDashboardStats godoc
// @Summary      Get dashboard statistics
// @Description  Returns overview stats (revenue, orders, products, customers)
// @Tags         Analytics
// @Security     BearerAuth
// @Success      200 {object} service.DashboardStats
// @Router       /admin/analytics/stats [get]
func (h *AnalyticsHandler) GetDashboardStats(c *gin.Context) {
	stats, err := h.Service.GetDashboardStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get stats"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": stats})
}

// GetSalesByCategory godoc
// @Summary      Get sales by category
// @Description  Returns revenue and count grouped by wine category
// @Tags         Analytics
// @Security     BearerAuth
// @Success      200 {array} service.SalesByCategory
// @Router       /admin/analytics/sales-by-category [get]
func (h *AnalyticsHandler) GetSalesByCategory(c *gin.Context) {
	data, err := h.Service.GetSalesByCategory()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

// GetTopProducts godoc
// @Summary      Get top selling products
// @Description  Returns top N best-selling wines
// @Tags         Analytics
// @Security     BearerAuth
// @Param        limit query int false "Number of products" default(5)
// @Success      200 {array} service.TopProduct
// @Router       /admin/analytics/top-products [get]
func (h *AnalyticsHandler) GetTopProducts(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))

	data, err := h.Service.GetTopProducts(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

// GetSalesByDay godoc
// @Summary      Get daily sales
// @Description  Returns sales data for the last N days
// @Tags         Analytics
// @Security     BearerAuth
// @Param        days query int false "Number of days" default(30)
// @Success      200 {array} service.SalesByDay
// @Router       /admin/analytics/sales-by-day [get]
func (h *AnalyticsHandler) GetSalesByDay(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))

	data, err := h.Service.GetSalesByDay(days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

// GetRecentOrders godoc
// @Summary      Get recent orders
// @Description  Returns the most recent orders
// @Tags         Analytics
// @Security     BearerAuth
// @Param        limit query int false "Number of orders" default(10)
// @Success      200 {array} service.RecentOrder
// @Router       /admin/analytics/recent-orders [get]
func (h *AnalyticsHandler) GetRecentOrders(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	data, err := h.Service.GetRecentOrders(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}
