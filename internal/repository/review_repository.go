package repository

import (
	"wine-shop-api/internal/domain"

	"gorm.io/gorm"
)

type reviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) ReviewRepository {
	return &reviewRepository{db: db}
}

func (r *reviewRepository) FindByProductAndUser(productID, userID uint) (*domain.Review, error) {
	var review domain.Review
	err := r.db.Where("product_id = ? AND user_id = ?", productID, userID).First(&review).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *reviewRepository) Create(review *domain.Review) error {
	return r.db.Create(review).Error
}

func (r *reviewRepository) PreloadUser(review *domain.Review) error {
	return r.db.Preload("User").First(review, review.ID).Error
}

func (r *reviewRepository) FindByProductID(productID uint) ([]domain.Review, error) {
	var reviews []domain.Review
	err := r.db.Preload("User").Where("product_id = ?", productID).
		Order("created_at desc").Find(&reviews).Error
	if err != nil {
		return nil, err
	}
	return reviews, nil
}

func (r *reviewRepository) AverageRating(productID uint) (float64, int64, error) {
	var result struct {
		Avg   float64
		Count int64
	}
	err := r.db.Model(&domain.Review{}).
		Select("COALESCE(AVG(rating), 0) as avg, COUNT(*) as count").
		Where("product_id = ?", productID).
		Scan(&result).Error
	if err != nil {
		return 0, 0, err
	}
	return result.Avg, result.Count, nil
}

func (r *reviewRepository) FindByID(id uint) (*domain.Review, error) {
	var review domain.Review
	if err := r.db.First(&review, id).Error; err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *reviewRepository) Delete(review *domain.Review) error {
	return r.db.Delete(review).Error
}
