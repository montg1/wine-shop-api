package service

import (
	"errors"

	"wine-shop-api/internal/domain"
	"wine-shop-api/pkg/config"
)

type ReviewService struct{}

// CreateReview creates a new review for a product
func (s *ReviewService) CreateReview(review *domain.Review) (*domain.Review, error) {
	// Check if user already reviewed this product
	var existing domain.Review
	if err := config.DB.Where("product_id = ? AND user_id = ?", review.ProductID, review.UserID).First(&existing).Error; err == nil {
		return nil, errors.New("you have already reviewed this product")
	}

	if err := config.DB.Create(review).Error; err != nil {
		return nil, err
	}

	// Load user data
	config.DB.Preload("User").First(review, review.ID)
	return review, nil
}

// GetProductReviews gets all reviews for a product
func (s *ReviewService) GetProductReviews(productID uint) ([]domain.Review, error) {
	var reviews []domain.Review
	if err := config.DB.Preload("User").Where("product_id = ?", productID).Order("created_at desc").Find(&reviews).Error; err != nil {
		return nil, err
	}
	return reviews, nil
}

// GetProductAverageRating calculates average rating for a product
func (s *ReviewService) GetProductAverageRating(productID uint) (float64, int64, error) {
	var result struct {
		Avg   float64
		Count int64
	}

	if err := config.DB.Model(&domain.Review{}).
		Select("COALESCE(AVG(rating), 0) as avg, COUNT(*) as count").
		Where("product_id = ?", productID).
		Scan(&result).Error; err != nil {
		return 0, 0, err
	}

	return result.Avg, result.Count, nil
}

// DeleteReview deletes a review (only by owner)
func (s *ReviewService) DeleteReview(reviewID, userID uint) error {
	var review domain.Review
	if err := config.DB.First(&review, reviewID).Error; err != nil {
		return errors.New("review not found")
	}

	if review.UserID != userID {
		return errors.New("you can only delete your own reviews")
	}

	return config.DB.Delete(&review).Error
}
