package application

import (
	"errors"
)

var (
	InvalidProductPriceErr = errors.New("Product Price Cannot be negative")
	InvalidProductQuantityErr = errors.New("Product Quantity cannot be negative")
	InvalidProductErr = errors.New("Product is invalid")
)

type Product struct {
	SKU string
	Name string
	Price float64
	Quantity int
}

func (p *Product) Validate() error {
	if p == nil {
		return InvalidProductErr
	}

	price := p.Price
	if price < 0 {
		return InvalidProductPriceErr
	}
	quantity := p.Quantity
	if quantity < 0 {
		return InvalidProductQuantityErr
	}
	return nil
}
