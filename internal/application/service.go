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

func (c *Checkout) Purchase(cartItems []domain.CartItem) (domain.Invoice, error) {
	invoiceItems, totalCost, err := c.applyPromotions(cartItems)
	if err != nil {
		return domain.Invoice{}, err
	}

	var reserveItems []domain.CartItem
	for _, invItem := range invoiceItems {
		reserveItems = append(reserveItems, domain.CartItem{
			SKU:      invItem.SKU,
			Quantity: invItem.Quantity,
		})
	}

	if err := c.purchaseRepo.BuyAllOrFail(reserveItems); err != nil {
		return domain.Invoice{}, err
	}

	return domain.Invoice{
		Items:     invoiceItems,
		TotalCost: totalCost,
	}, nil
}

func (c *Checkout) applyPromotions(items []domain.CartItem) ([]domain.InvoiceItem, float64, error) {
	var invoiceItems []domain.InvoiceItem
	var totalCost float64

	for _, item := range items {
		product, err := c.purchaseRepo.GetProduct(item.SKU)
		if err != nil {
			return nil, 0, err
		}

		switch item.SKU {
		case "43N23P":
			macInvoice := domain.InvoiceItem{
				SKU:        product.SKU,
				Name:       product.Name,
				Quantity:   item.Quantity,
				UnitPrice:  product.Price,
				TotalPrice: product.Price * float64(item.Quantity),
			}
			invoiceItems = append(invoiceItems, macInvoice)
			totalCost += macInvoice.TotalPrice

			// Free Raspberry Pi B line (price is $0).
			piProduct, err := c.purchaseRepo.GetProduct("234234")
			if err == nil {
				freePi := domain.InvoiceItem{
					SKU:        piProduct.SKU,
					Name:       piProduct.Name,
					Quantity:   item.Quantity, // one free per MacBook Pro
					UnitPrice:  0.0,
					TotalPrice: 0.0,
				}
				invoiceItems = append(invoiceItems, freePi)
			}
		case "120P90":
			// Google TV: Buy 3 for the price of 2.
			chargedQty := (item.Quantity/3)*2 + (item.Quantity % 3)
			invoiceItem := domain.InvoiceItem{
				SKU:        product.SKU,
				Name:       product.Name,
				Quantity:   item.Quantity,
				UnitPrice:  product.Price,
				TotalPrice: product.Price * float64(chargedQty),
			}
			invoiceItems = append(invoiceItems, invoiceItem)
			totalCost += invoiceItem.TotalPrice
		case "A304SD":
			// Alexa Speaker: 10% discount if buying more than 3.
			var unitPrice float64
			if item.Quantity > 3 {
				unitPrice = product.Price * 0.9
			} else {
				unitPrice = product.Price
			}
			invoiceItem := domain.InvoiceItem{
				SKU:        product.SKU,
				Name:       product.Name,
				Quantity:   item.Quantity,
				UnitPrice:  unitPrice,
				TotalPrice: unitPrice * float64(item.Quantity),
			}
			invoiceItems = append(invoiceItems, invoiceItem)
			totalCost += invoiceItem.TotalPrice
		default:
			// No promotion.
			invoiceItem := domain.InvoiceItem{
				SKU:        product.SKU,
				Name:       product.Name,
				Quantity:   item.Quantity,
				UnitPrice:  product.Price,
				TotalPrice: product.Price * float64(item.Quantity),
			}
			invoiceItems = append(invoiceItems, invoiceItem)
			totalCost += invoiceItem.TotalPrice
		}
	}

	return invoiceItems, totalCost, nil
}
