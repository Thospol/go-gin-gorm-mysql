package product

import (
	"go-gin-gorm-mysql/internal/core/config"
	"go-gin-gorm-mysql/internal/pkg/models"
	"go-gin-gorm-mysql/internal/pkg/repository"

	"gorm.io/gorm"
)

// Service service product interface
type Service interface {
	Create(database *gorm.DB, request createRequest) (*models.Product, error)
	GetAll(database *gorm.DB) ([]*models.Product, error)
	GetByID(database *gorm.DB, id uint) (*models.Product, error)
}

type service struct {
	config            *config.Configs
	result            *config.ReturnResult
	productRepository repository.Product
}

// NewService new service product
func NewService(config *config.Configs, result *config.ReturnResult) Service {
	return &service{
		config:            config,
		result:            result,
		productRepository: repository.NewProductRepository(),
	}
}

// Create create product service
func (s *service) Create(database *gorm.DB, request createRequest) (*models.Product, error) {
	product := &models.Product{
		Name:        request.Name,
		Description: request.Name,
		Price:       request.Price,
		Amount:      request.Amount,
	}
	err := s.productRepository.Create(database, product)
	if err != nil {
		return nil, err
	}

	return product, nil
}

// GetAll get all product service
func (s *service) GetAll(database *gorm.DB) ([]*models.Product, error) {
	entities, err := s.productRepository.FindAll(database)
	if err != nil {
		return nil, err
	}

	return entities, nil
}

// GetByID get product by id service
func (s *service) GetByID(database *gorm.DB, id uint) (*models.Product, error) {
	product := &models.Product{}
	err := s.productRepository.FindByID(database, id, product)
	if err != nil {
		return nil, s.result.Internal.DatabaseNotFound
	}

	return product, nil
}
