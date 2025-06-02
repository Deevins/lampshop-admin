package repository

import (
	"github.com/Deevins/lampshop-admin-backend/internal/domain"
)

// CategoryRepository описывает методы для работы с категориями.
type CategoryRepository interface {
	GetAll() ([]domain.Category, error)
	GetByID(id string) (*domain.Category, error)
}

// InMemoryCategoryRepo — простая in-memory реализация CategoryRepository.
type InMemoryCategoryRepo struct {
	categories []domain.Category
}

// NewInMemoryCategoryRepo создаёт новый репозиторий с начальными категориями.
func NewInMemoryCategoryRepo(initial []domain.Category) *InMemoryCategoryRepo {
	return &InMemoryCategoryRepo{categories: initial}
}

func (r *InMemoryCategoryRepo) GetAll() ([]domain.Category, error) {
	cats := make([]domain.Category, len(r.categories))
	copy(cats, r.categories)
	return cats, nil
}

func (r *InMemoryCategoryRepo) GetByID(id string) (*domain.Category, error) {
	for _, c := range r.categories {
		if c.ID == id {
			copyCat := c
			return &copyCat, nil
		}
	}
	return nil, ErrCategoryNotFound
}

// RepoError — простая обёртка для ошибок репозитория.
type RepoError struct {
	Message string
}

func (e *RepoError) Error() string {
	return e.Message
}

var ErrCategoryNotFound = &RepoError{"category not found"}
