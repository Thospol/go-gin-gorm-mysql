package repository

import (
	"go-gin-gorm-mysql/internal/pkg/models"

	"gorm.io/gorm"
)

// Product product repository
type Product interface {
	Create(database *gorm.DB, i interface{}) error
	FindByName(database *gorm.DB, name string) ([]*models.Product, error)
}

type product struct {
	Repository
}

// NewProductRepository new product repository
func NewProductRepository() Product {
	return &product{
		NewRepository(),
	}
}

// FindByName find product where name follow condition.
func (repo *product) FindByName(database *gorm.DB, name string) ([]*models.Product, error) {
	entities := []*models.Product{}
	if err := database.Where("name = ?", name).Find(&entities).Error; err != nil {
		return nil, err
	}

	return entities, nil
}
