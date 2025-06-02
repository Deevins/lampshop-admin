package service

import (
	"github.com/Deevins/lampshop-admin-backend/internal/domain"
	"github.com/Deevins/lampshop-admin-backend/internal/repository"
)

// OrderService описывает бизнес-логику для заказов.
type OrderService interface {
	GetAllOrders() ([]domain.Order, error)
	UpdateOrderStatus(id int, status domain.OrderStatus) (*domain.Order, error)
}

type orderService struct {
	repo repository.OrderRepository
}

// NewOrderService создаёт новый OrderService.
func NewOrderService(r repository.OrderRepository) OrderService {
	return &orderService{repo: r}
}

func (s *orderService) GetAllOrders() ([]domain.Order, error) {
	return s.repo.GetAll()
}

func (s *orderService) UpdateOrderStatus(id int, status domain.OrderStatus) (*domain.Order, error) {
	return s.repo.UpdateStatus(id, status)
}
