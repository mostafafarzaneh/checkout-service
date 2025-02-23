package port

import (
	"checkout/internal/domain"
)

type PurchaseRepository interface {
	BuyAllOrFail(cart []domain.CartItem) error

	GetProduct(sku string) (domain.Product, error)
	UpdateProduct(domain.Product) error
}
