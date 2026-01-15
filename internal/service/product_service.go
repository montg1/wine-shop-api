package service

import (
	"errors"

	"wine-shop-api/internal/domain"
	"wine-shop-api/pkg/config"
)

type ProductService struct{}

func (s *ProductService) CreateProduct(product *domain.Product) (*domain.Product, error) {
	if err := config.DB.Create(&product).Error; err != nil {
		return nil, err
	}
	return product, nil
}

func (s *ProductService) GetAllProducts(page, limit int, search, category string) ([]domain.Product, int64, error) {
	var products []domain.Product
	var total int64

	offset := (page - 1) * limit

	query := config.DB.Model(&domain.Product{})

	// Apply search filter
	if search != "" {
		query = query.Where("LOWER(name) LIKE LOWER(?)", "%"+search+"%")
	}

	// Apply category filter
	if category != "" {
		query = query.Where("LOWER(category) = LOWER(?)", category)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(limit).Offset(offset).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (s *ProductService) GetProductByID(id uint) (*domain.Product, error) {
	var product domain.Product
	if err := config.DB.First(&product, id).Error; err != nil {
		return nil, errors.New("product not found")
	}
	return &product, nil
}

func (s *ProductService) UpdateProduct(id uint, input *domain.Product) (*domain.Product, error) {
	var product domain.Product
	if err := config.DB.First(&product, id).Error; err != nil {
		return nil, errors.New("product not found")
	}

	// Update fields
	product.Name = input.Name
	product.Description = input.Description
	product.Price = input.Price
	product.Stock = input.Stock
	product.ImageURL = input.ImageURL
	product.Category = input.Category

	if err := config.DB.Save(&product).Error; err != nil {
		return nil, err
	}

	return &product, nil
}

func (s *ProductService) DeleteProduct(id uint) error {
	if err := config.DB.Delete(&domain.Product{}, id).Error; err != nil {
		return err
	}
	return nil
}
