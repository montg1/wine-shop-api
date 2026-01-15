package domain

import "gorm.io/gorm"

type Review struct {
	gorm.Model
	ProductID uint    `json:"product_id"`
	UserID    uint    `json:"user_id"`
	Rating    int     `json:"rating" binding:"required,min=1,max=5"`
	Comment   string  `json:"comment"`
	User      User    `json:"user" gorm:"foreignKey:UserID"`
	Product   Product `json:"-" gorm:"foreignKey:ProductID"`
}
