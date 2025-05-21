package models

import (
	"gorm.io/gorm"
)

// Book 定义图书模型
type Book struct {
	gorm.Model
	Title       string `json:"title" gorm:"size:100;not null"`
	Author      string `json:"author" gorm:"size:100;not null"`
	Description string `json:"description" gorm:"type:text"`
	Price       float64 `json:"price" gorm:"type:decimal(10,2)"`
	Stock       int     `json:"stock" gorm:"not null;default:0"`
}