package port

import (
	"checkout/internal/domain"
)

type OrderProcessor interface {
	ProcessOrder(order []domain.OrderItem) error
}
