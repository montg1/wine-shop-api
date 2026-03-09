package service

import (
	"errors"
	"html"
	"strings"

	"wine-shop-api/internal/domain"
	"wine-shop-api/internal/repository"
)

type ReviewService struct {
	Repo repository.ReviewRepository
}

func NewReviewService(repo repository.ReviewRepository) *ReviewService {
	return &ReviewService{Repo: repo}
}

// CreateReview creates a new review for a product
func (s *ReviewService) CreateReview(review *domain.Review) (*domain.Review, error) {
	// Check if user already reviewed this product
	if _, err := s.Repo.FindByProductAndUser(review.ProductID, review.UserID); err == nil {
		return nil, errors.New("you have already reviewed this product")
	}

	// Sanitize user-generated content
	review.Comment = html.EscapeString(strings.TrimSpace(review.Comment))

	if err := s.Repo.Create(review); err != nil {
		return nil, err
	}

	// Load user data
	s.Repo.PreloadUser(review)
	return review, nil
}

// GetProductReviews gets all reviews for a product
func (s *ReviewService) GetProductReviews(productID uint) ([]domain.Review, error) {
	return s.Repo.FindByProductID(productID)
}

// GetProductAverageRating calculates average rating for a product
func (s *ReviewService) GetProductAverageRating(productID uint) (float64, int64, error) {
	return s.Repo.AverageRating(productID)
}

// DeleteReview deletes a review (only by owner)
func (s *ReviewService) DeleteReview(reviewID, userID uint) error {
	review, err := s.Repo.FindByID(reviewID)
	if err != nil {
		return errors.New("review not found")
	}

	if review.UserID != userID {
		return errors.New("you can only delete your own reviews")
	}

	return s.Repo.Delete(review)
}
