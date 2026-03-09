package service

import (
	"errors"
	"html"
	"strings"

	"wine-shop-api/internal/domain"
	"wine-shop-api/internal/repository"
)

type ProductService struct {
	Repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) *ProductService {
	return &ProductService{Repo: repo}
}

func (s *ProductService) CreateProduct(product *domain.Product) (*domain.Product, error) {
	// Sanitize user-generated content
	product.Name = html.EscapeString(strings.TrimSpace(product.Name))
	product.Description = html.EscapeString(strings.TrimSpace(product.Description))

	if err := s.Repo.Create(product); err != nil {
		return nil, err
	}
	return product, nil
}

func (s *ProductService) GetAllProducts(page, limit int, search, category string) ([]domain.Product, int64, error) {
	offset := (page - 1) * limit
	return s.Repo.FindAll(offset, limit, search, category)
}

func (s *ProductService) GetProductByID(id uint) (*domain.Product, error) {
	product, err := s.Repo.FindByID(id)
	if err != nil {
		return nil, errors.New("product not found")
	}
	return product, nil
}

func (s *ProductService) UpdateProduct(id uint, input *domain.Product) (*domain.Product, error) {
	product, err := s.Repo.FindByID(id)
	if err != nil {
		return nil, errors.New("product not found")
	}

	// Sanitize user-generated content
	product.Name = html.EscapeString(strings.TrimSpace(input.Name))
	product.Description = html.EscapeString(strings.TrimSpace(input.Description))
	product.Price = input.Price
	product.Stock = input.Stock
	product.ImageURL = input.ImageURL
	product.Category = input.Category

	if err := s.Repo.Save(product); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *ProductService) DeleteProduct(id uint) error {
	return s.Repo.Delete(id)
}
