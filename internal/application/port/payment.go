package port

import (
	"checkout/internal/domain"
)

type PaymentProcessor interface {
	ProcessPayment(order []domain.OrderItem) error
}
