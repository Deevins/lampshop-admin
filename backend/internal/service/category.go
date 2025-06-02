package service

import (
	"github.com/Deevins/lampshop-admin-backend/internal/domain"
	"github.com/Deevins/lampshop-admin-backend/internal/repository"
)

// CategoryService описывает бизнес-логику по категориям.
type CategoryService interface {
	GetAllCategories() ([]domain.Category, error)
	GetAttributeOptions(categoryID string) ([]domain.AttributeOption, error)
}

type categoryService struct {
	catRepo  repository.CategoryRepository
	attrRepo repository.AttributeRepository
}

// NewCategoryService создаёт новый CategoryService.
func NewCategoryService(
	cr repository.CategoryRepository,
	ar repository.AttributeRepository,
) CategoryService {
	return &categoryService{catRepo: cr, attrRepo: ar}
}

func (s *categoryService) GetAllCategories() ([]domain.Category, error) {
	return s.catRepo.GetAll()
}

func (s *categoryService) GetAttributeOptions(categoryID string) ([]domain.AttributeOption, error) {
	return s.attrRepo.GetByCategoryID(categoryID)
}
