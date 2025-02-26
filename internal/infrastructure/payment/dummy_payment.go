package payment

import (
	"checkout/internal/domain"
)

type DummyPaymentProcessor struct{}

func (p *DummyPaymentProcessor) ProcessPayment(order []domain.OrderItem) error {
	return nil
}
