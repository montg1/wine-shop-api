package service

import (
	"errors"

	"wine-shop-api/internal/domain"
	"wine-shop-api/internal/repository"

	"gorm.io/gorm"
)

type CartService struct {
	Repo        repository.CartRepository
	ProductRepo repository.ProductRepository
}

func NewCartService(repo repository.CartRepository, productRepo repository.ProductRepository) *CartService {
	return &CartService{Repo: repo, ProductRepo: productRepo}
}

func (s *CartService) GetCart(userID uint) (*domain.Cart, error) {
	cart, err := s.Repo.FindByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create a new cart if one doesn't exist
			cart = &domain.Cart{UserID: userID}
			if err := s.Repo.Create(cart); err != nil {
				return nil, err
			}
			return cart, nil
		}
		return nil, err
	}
	return cart, nil
}

func (s *CartService) AddToCart(userID uint, productID uint, quantity int) error {
	cart, err := s.GetCart(userID)
	if err != nil {
		return err
	}

	// Check if product exists
	if _, err := s.ProductRepo.FindByID(productID); err != nil {
		return errors.New("product not found")
	}

	// Check if item already exists in cart
	cartItem, err := s.Repo.FindCartItem(cart.ID, productID)
	if err == nil {
		// Update quantity
		cartItem.Quantity += quantity
		return s.Repo.SaveCartItem(cartItem)
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		// Create new item
		newItem := &domain.CartItem{
			CartID:    cart.ID,
			ProductID: productID,
			Quantity:  quantity,
		}
		return s.Repo.CreateCartItem(newItem)
	}

	return err
}

func (s *CartService) ClearCart(userID uint) error {
	cart, err := s.GetCart(userID)
	if err != nil {
		return err
	}
	return s.Repo.DeleteCartItems(cart.ID)
}
