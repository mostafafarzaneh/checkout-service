package application


type CartItem struct {
	SKU string
	Quantity int
}

type Checkout struct {
	Repository *MemoryDB
}

func NewCheckout(Repository *MemoryDB) *Checkout {
	return &Checkout{Repository}
}

func (c *Checkout) Purchase(cartItems []CartItem) (float64, error) {
	var totalCost float64 = 0
	for _, item := range cartItems {
		p, err := c.Repository.GetProduct(item.SKU)
		if err != nil {
			return -1, err
		}

		if err := c.Repository.BuyProduct(item.SKU, item.Quantity); err != nil {
			return -1, err
		}
		totalCost += p.Price * float64(item.Quantity)
	}
	return totalCost, nil
}
