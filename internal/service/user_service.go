package service

import (
	"errors"
	"html"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"wine-shop-api/internal/domain"
	"wine-shop-api/internal/repository"
	"wine-shop-api/pkg/utils"
)

type UserService struct {
	Repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{Repo: repo}
}

func (s *UserService) Register(user *domain.User) (*domain.User, error) {
	// 1. Check if email exists
	if _, err := s.Repo.FindByEmail(user.Email); err == nil {
		return nil, errors.New("email already in use")
	}

	// 2. Hash Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashedPassword)
	user.Email = html.EscapeString(strings.TrimSpace(user.Email))

	// 3. Create User
	if err := s.Repo.Create(user); err != nil {
		return nil, err
	}

	// Remove password from response
	user.Password = ""
	return user, nil
}

func (s *UserService) Login(email, password string) (string, error) {
	// 1. Find User
	user, err := s.Repo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("invalid email or password")
		}
		return "", err
	}

	// 2. Verify Password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", errors.New("invalid email or password")
	}

	// 3. Generate Token
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

// PromoteToAdmin promotes a user to admin role
func (s *UserService) PromoteToAdmin(userID uint) error {
	user, err := s.Repo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}
	user.Role = "admin"
	return s.Repo.Save(user)
}

// GetUserByID returns a user by ID
func (s *UserService) GetUserByID(userID uint) (*domain.User, error) {
	user, err := s.Repo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}
