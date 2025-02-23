package application


type Checkout struct {
	Repository *MemoryDB
}

func NewCheckout(Repository *MemoryDB) *Checkout {
	return &Checkout{Repository}
}

func (c *Checkout) Purchase(cartItems []CartItem) (float64, error) {
	return c.Repository.BuyProducts(cartItems)
}
