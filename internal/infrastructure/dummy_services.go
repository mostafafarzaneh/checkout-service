package infrastructure

import (
	"checkout/internal/domain"
)

type DummyPaymentProcessor struct{}

func (p *DummyPaymentProcessor) ProcessPayment(order []domain.OrderItem) error {
	return nil
}

type DummyOrderProcessor struct{}

func (p *DummyOrderProcessor) ProcessOrder(order []domain.OrderItem) error {
	return nil
}
