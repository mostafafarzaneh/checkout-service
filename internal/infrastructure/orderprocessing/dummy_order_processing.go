package orderprocessing

import (
	"checkout/internal/domain"
)

type DummyOrderProcessor struct{}

func (p *DummyOrderProcessor) ProcessOrder(order []domain.OrderItem) error {
	return nil
}
