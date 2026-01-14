package service

import (
	"errors"
	"html"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"wine-shop-api/internal/domain"
	"wine-shop-api/pkg/config"
	"wine-shop-api/pkg/utils"
)

type UserService struct{}

func (s *UserService) Register(user *domain.User) (*domain.User, error) {
	// 1. Check if email exists
	var existingUser domain.User
	if err := config.DB.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
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
	if err := config.DB.Create(&user).Error; err != nil {
		return nil, err
	}

	// Remove password from response
	user.Password = ""
	return user, nil
}

func (s *UserService) Login(email, password string) (string, error) {
	var user domain.User

	// 1. Find User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("invalid email or password")
		}
		return "", err
	}

	// 2. Verify Password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
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
