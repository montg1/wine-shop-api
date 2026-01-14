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

// CreateOrder godoc
// @Summary      Checkout (Place Order)
// @Description  Convert current cart into an order and clear the cart
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      201    {object}  map[string]interface{}
// @Failure      400    {object}  map[string]interface{}
// @Failure      401    {object}  map[string]interface{}
// @Router       /orders [post]
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

// GetOrders godoc
// @Summary      Get order history
// @Description  List past orders for the authenticated user
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200    {object}  map[string]interface{}
// @Failure      401    {object}  map[string]interface{}
// @Failure      500    {object}  map[string]interface{}
// @Router       /orders [get]
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
