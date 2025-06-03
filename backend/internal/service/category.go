package service

import (
	"context"
	"github.com/Deevins/lampshop-admin-backend/internal/repository"
	"github.com/Deevins/lampshop-admin-backend/internal/repository/sql"
	"github.com/google/uuid"
)

type CategoryService interface {
	GetAllCategories(ctx context.Context) ([]*sql.Category, error)
	GetAttributeOptions(ctx context.Context, categoryID uuid.UUID) ([]*sql.AttributeOption, error)
}

type categoryService struct {
	repo repository.AttributeRepository
}

func NewCategoryService(repo repository.AttributeRepository) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) GetAllCategories(ctx context.Context) ([]*sql.Category, error) {
	return s.repo.ListCategories(ctx)
}

func (s *categoryService) GetAttributeOptions(ctx context.Context, categoryID uuid.UUID) ([]*sql.AttributeOption, error) {
	return s.repo.ListAttributesByCategory(ctx, categoryID)
}
