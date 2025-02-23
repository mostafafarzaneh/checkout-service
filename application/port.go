package application

type PurchaseRepository interface {
	BuyAllOrFail(cart []CartItem) error

	GetProduct(sku string) (Product, error)
	UpdateProduct(Product) error
}
