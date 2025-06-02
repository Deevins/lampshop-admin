package repository

import "github.com/Deevins/lampshop-admin-backend/internal/domain"

// AttributeRepository описывает методы для получения атрибутов по категории.
type AttributeRepository interface {
	GetByCategoryID(categoryID string) ([]domain.AttributeOption, error)
}

// InMemoryAttributeRepo — in-memory реализация AttributeRepository.
type InMemoryAttributeRepo struct {
	options map[string][]domain.AttributeOption
}

// NewInMemoryAttributeRepo создаёт репозиторий с начальными опциями.
func NewInMemoryAttributeRepo(initial map[string][]domain.AttributeOption) *InMemoryAttributeRepo {
	return &InMemoryAttributeRepo{options: initial}
}

func (r *InMemoryAttributeRepo) GetByCategoryID(categoryID string) ([]domain.AttributeOption, error) {
	if opts, ok := r.options[categoryID]; ok {
		copyOpts := make([]domain.AttributeOption, len(opts))
		copy(copyOpts, opts)
		return copyOpts, nil
	}
	return nil, ErrAttributesNotFound
}

var ErrAttributesNotFound = &RepoError{"attributes not found for category"}
