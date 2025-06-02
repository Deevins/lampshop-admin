package repository

import (
	"errors"
	"github.com/Deevins/lampshop-admin-backend/internal/domain"
	"sync"
	"time"
)

// ProductRepository описывает методы CRUD для товаров.
type ProductRepository interface {
	GetAll() ([]domain.Product, error)
	GetByID(id int) (*domain.Product, error)
	Create(prod domain.Product) (*domain.Product, error)
	Update(id int, prod domain.Product) (*domain.Product, error)
	Delete(id int) error
}

// InMemoryProductRepo — in-memory реализация ProductRepository.
type InMemoryProductRepo struct {
	mu       sync.RWMutex
	products []domain.Product
}

// NewInMemoryProductRepo создаёт новый репозиторий с начальными товарами.
func NewInMemoryProductRepo(initial []domain.Product) *InMemoryProductRepo {
	return &InMemoryProductRepo{products: initial}
}

func (r *InMemoryProductRepo) GetAll() ([]domain.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := make([]domain.Product, len(r.products))
	copy(list, r.products)
	return list, nil
}

func (r *InMemoryProductRepo) GetByID(id int) (*domain.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, p := range r.products {
		if p.ID == id {
			copyProd := p
			return &copyProd, nil
		}
	}
	return nil, ErrProductNotFound
}

func (r *InMemoryProductRepo) Create(prod domain.Product) (*domain.Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	maxID := 0
	for _, p := range r.products {
		if p.ID > maxID {
			maxID = p.ID
		}
	}
	prod.ID = maxID + 1
	prod.CreatedAt = time.Now()
	prod.UpdatedAt = time.Now()
	r.products = append(r.products, prod)

	copyProd := prod
	return &copyProd, nil
}

func (r *InMemoryProductRepo) Update(id int, prod domain.Product) (*domain.Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, p := range r.products {
		if p.ID == id {
			prod.ID = id
			prod.CreatedAt = p.CreatedAt
			prod.UpdatedAt = time.Now()
			r.products[i] = prod

			copyProd := prod
			return &copyProd, nil
		}
	}
	return nil, ErrProductNotFound
}

func (r *InMemoryProductRepo) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, p := range r.products {
		if p.ID == id {
			r.products = append(r.products[:i], r.products[i+1:]...)
			return nil
		}
	}
	return ErrProductNotFound
}

var ErrProductNotFound = errors.New("product not found")
