package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"wine-shop-api/internal/domain"
	"wine-shop-api/internal/service"
	"wine-shop-api/pkg/utils"
)

type ReviewHandler struct {
	Service *service.ReviewService
}

// CreateReview godoc
// @Summary      Create a review
// @Description  Add a review and rating for a product
// @Tags         Reviews
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id     path      int            true  "Product ID"
// @Param        input  body      domain.Review  true  "Review Data"
// @Success      201    {object}  domain.Review
// @Failure      400    {object}  map[string]interface{}
// @Router       /products/{id}/reviews [post]
func (h *ReviewHandler) CreateReview(c *gin.Context) {
	productID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	userID, err := utils.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input struct {
		Rating  int    `json:"rating" binding:"required,min=1,max=5"`
		Comment string `json:"comment"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	review := &domain.Review{
		ProductID: uint(productID),
		UserID:    userID,
		Rating:    input.Rating,
		Comment:   input.Comment,
	}

	createdReview, err := h.Service.CreateReview(review)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": createdReview})
}

// GetProductReviews godoc
// @Summary      Get product reviews
// @Description  Get all reviews for a specific product
// @Tags         Reviews
// @Accept       json
// @Produce      json
// @Param        id  path      int  true  "Product ID"
// @Success      200 {object}  map[string]interface{}
// @Router       /products/{id}/reviews [get]
func (h *ReviewHandler) GetProductReviews(c *gin.Context) {
	productID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	reviews, err := h.Service.GetProductReviews(uint(productID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	avgRating, count, _ := h.Service.GetProductAverageRating(uint(productID))

	c.JSON(http.StatusOK, gin.H{
		"data": reviews,
		"meta": gin.H{
			"average_rating": avgRating,
			"total_reviews":  count,
		},
	})
}

// DeleteReview godoc
// @Summary      Delete a review
// @Description  Delete your own review
// @Tags         Reviews
// @Security     BearerAuth
// @Param        id        path      int  true  "Product ID"
// @Param        reviewId  path      int  true  "Review ID"
// @Success      200       {object}  map[string]interface{}
// @Failure      400       {object}  map[string]interface{}
// @Router       /products/{id}/reviews/{reviewId} [delete]
func (h *ReviewHandler) DeleteReview(c *gin.Context) {
	reviewID, err := strconv.Atoi(c.Param("reviewId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	userID, err := utils.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := h.Service.DeleteReview(uint(reviewID), userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Review deleted successfully"})
}
