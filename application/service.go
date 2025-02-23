package application


type Checkout struct {
	purchaseRepo PurchaseRepository
}

func NewCheckout(repo PurchaseRepository) *Checkout {
	return &Checkout{repo}
}

func (c *Checkout) Purchase(cartItems []CartItem) (float64, error) {
	return c.purchaseRepo.BuyAllOrFail(cartItems)
}
