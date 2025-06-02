package service

import (
	"github.com/Deevins/lampshop-admin-backend/internal/domain"
	"github.com/Deevins/lampshop-admin-backend/internal/repository"
)

// ProductService описывает бизнес-логику для товаров.
type ProductService interface {
	GetAllProducts() ([]domain.Product, error)
	GetProductByID(id int) (*domain.Product, error)
	CreateProduct(input domain.Product) (*domain.Product, error)
	UpdateProduct(id int, input domain.Product) (*domain.Product, error)
	DeleteProduct(id int) error
}

type productService struct {
	repo repository.ProductRepository
}

// NewProductService создаёт новый ProductService.
func NewProductService(r repository.ProductRepository) ProductService {
	return &productService{repo: r}
}

func (s *productService) GetAllProducts() ([]domain.Product, error) {
	return s.repo.GetAll()
}

func (s *productService) GetProductByID(id int) (*domain.Product, error) {
	return s.repo.GetByID(id)
}

func (s *productService) CreateProduct(input domain.Product) (*domain.Product, error) {
	return s.repo.Create(input)
}

func (s *productService) UpdateProduct(id int, input domain.Product) (*domain.Product, error) {
	return s.repo.Update(id, input)
}

func (s *productService) DeleteProduct(id int) error {
	return s.repo.Delete(id)
}
