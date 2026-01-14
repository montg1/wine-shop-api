package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"wine-shop-api/internal/service"
	"wine-shop-api/pkg/utils"
)

type OrderHandler struct {
	Service *service.OrderService
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	userID, err := utils.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	order, err := h.Service.CreateOrder(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Order placed successfully", "data": order})
}

func (h *OrderHandler) GetOrders(c *gin.Context) {
	userID, err := utils.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	orders, err := h.Service.GetOrders(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": orders})
}
