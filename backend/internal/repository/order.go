package repository

import (
	"errors"
	"github.com/Deevins/lampshop-admin-backend/internal/domain"
	"sync"
	"time"
)

// OrderRepository описывает методы работы с заказами.
type OrderRepository interface {
	GetAll() ([]domain.Order, error)
	GetByID(id int) (*domain.Order, error)
	UpdateStatus(id int, status domain.OrderStatus) (*domain.Order, error)
}

// InMemoryOrderRepo — in-memory реализация OrderRepository.
type InMemoryOrderRepo struct {
	mu     sync.RWMutex
	orders []domain.Order
}

// NewInMemoryOrderRepo создаёт репозиторий с начальными заказами.
func NewInMemoryOrderRepo(initial []domain.Order) *InMemoryOrderRepo {
	return &InMemoryOrderRepo{orders: initial}
}

func (r *InMemoryOrderRepo) GetAll() ([]domain.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := make([]domain.Order, len(r.orders))
	copy(list, r.orders)
	return list, nil
}

func (r *InMemoryOrderRepo) GetByID(id int) (*domain.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, o := range r.orders {
		if o.ID == id {
			copyOrder := o
			return &copyOrder, nil
		}
	}
	return nil, ErrOrderNotFound
}

func (r *InMemoryOrderRepo) UpdateStatus(id int, status domain.OrderStatus) (*domain.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, o := range r.orders {
		if o.ID == id {
			r.orders[i].Status = status
			r.orders[i].UpdatedAt = time.Now()

			copyOrder := r.orders[i]
			return &copyOrder, nil
		}
	}
	return nil, ErrOrderNotFound
}

var ErrOrderNotFound = errors.New("order not found")
