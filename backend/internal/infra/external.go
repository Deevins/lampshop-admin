package infra

import (
	"errors"
)

// ExternalClient описывает интерфейс для вызовов внешних сервисов.
// Пример: отправка уведомления, логирование изменения статуса заказа и т. д.
type ExternalClient interface {
	NotifyOrderStatusChange(orderID int, newStatus string) error
}

// InMemoryExternalClient — “заглушка” для ExternalClient,
// пока действительно не нужно делать реальный HTTP-вызов.
type InMemoryExternalClient struct{}

func NewInMemoryExternalClient() *InMemoryExternalClient {
	return &InMemoryExternalClient{}
}

func (c *InMemoryExternalClient) NotifyOrderStatusChange(orderID int, newStatus string) error {
	// Здесь мог бы быть HTTP POST на внешний API и т. д.
	// Пока просто “заглушка” — ничего не делаем и возвращаем nil.
	return nil
}

// RealExternalClient — пример, как мог бы выглядеть реальный клиент.
// На данный момент он просто возвращает “не реализовано”.
type RealExternalClient struct {
	// Например, вы могли бы хранить HTTP-Client, endpoint и т. д.
}

func NewRealExternalClient() *RealExternalClient {
	return &RealExternalClient{}
}

func (c *RealExternalClient) NotifyOrderStatusChange(orderID int, newStatus string) error {
	// Здесь вы бы сделали HTTP-запрос, токен OAuth, логирование и т.д.
	return errors.New("not implemented")
}
