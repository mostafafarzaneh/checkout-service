package application

import (
	"checkout/internal/application/port"
	"checkout/internal/domain"
)

type Checkout struct {
	purchaseRepo port.PurchaseRepository
}

func NewCheckout(repo port.PurchaseRepository) *Checkout {
	return &Checkout{repo}
}

func (c *Checkout) Purchase(cartItems []domain.CartItem) (domain.Order, error) {
	orderItems, totalCost, err := c.applyPromotions(cartItems)
	if err != nil {
		return domain.Order{}, err
	}

	var reserveItems []domain.CartItem
	for _, invItem := range orderItems {
		reserveItems = append(reserveItems, domain.CartItem{
			SKU:      invItem.SKU,
			Quantity: invItem.Quantity,
		})
	}

	if err := c.purchaseRepo.BuyAllOrFail(reserveItems); err != nil {
		return domain.Order{}, err
	}

	return domain.Order{
		Items:     orderItems,
		TotalCost: totalCost,
	}, nil
}

func (c *Checkout) applyPromotions(items []domain.CartItem) ([]domain.OrderItem, float64, error) {
	var orderItems []domain.OrderItem
	var totalCost float64

	for _, item := range items {
		product, err := c.purchaseRepo.GetProduct(item.SKU)
		if err != nil {
			return nil, 0, err
		}

		switch item.SKU {
		case "43N23P":
			macOrder := domain.OrderItem{
				SKU:        product.SKU,
				Name:       product.Name,
				Quantity:   item.Quantity,
				UnitPrice:  product.Price,
				TotalPrice: product.Price * float64(item.Quantity),
			}
			orderItems = append(orderItems, macOrder)
			totalCost += macOrder.TotalPrice

			// Free Raspberry Pi B line (price is $0).
			piProduct, err := c.purchaseRepo.GetProduct("234234")
			if err == nil {
				freePi := domain.OrderItem{
					SKU:        piProduct.SKU,
					Name:       piProduct.Name,
					Quantity:   item.Quantity, // one free per MacBook Pro
					UnitPrice:  0.0,
					TotalPrice: 0.0,
				}
				orderItems = append(orderItems, freePi)
			}
		case "120P90":
			// Google TV: Buy 3 for the price of 2.
			chargedQty := (item.Quantity/3)*2 + (item.Quantity % 3)
			OrderItem := domain.OrderItem{
				SKU:        product.SKU,
				Name:       product.Name,
				Quantity:   item.Quantity,
				UnitPrice:  product.Price,
				TotalPrice: product.Price * float64(chargedQty),
			}
			orderItems = append(orderItems, OrderItem)
			totalCost += OrderItem.TotalPrice
		case "A304SD":
			// Alexa Speaker: 10% discount if buying more than 3.
			var unitPrice float64
			if item.Quantity > 3 {
				unitPrice = product.Price * 0.9
			} else {
				unitPrice = product.Price
			}
			OrderItem := domain.OrderItem{
				SKU:        product.SKU,
				Name:       product.Name,
				Quantity:   item.Quantity,
				UnitPrice:  unitPrice,
				TotalPrice: unitPrice * float64(item.Quantity),
			}
			orderItems = append(orderItems, OrderItem)
			totalCost += OrderItem.TotalPrice
		default:
			// No promotion.
			OrderItem := domain.OrderItem{
				SKU:        product.SKU,
				Name:       product.Name,
				Quantity:   item.Quantity,
				UnitPrice:  product.Price,
				TotalPrice: product.Price * float64(item.Quantity),
			}
			orderItems = append(orderItems, OrderItem)
			totalCost += OrderItem.TotalPrice
		}
	}

	return orderItems, totalCost, nil
}
