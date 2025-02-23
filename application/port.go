package application

type PurchaseRepository interface {
	BuyAllOrFail(cart []CartItem) (float64, error)

	GetProduct(sku string) (Product, error)
	UpdateProduct(Product) error
}
