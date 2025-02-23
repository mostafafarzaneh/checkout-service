package application


type Checkout struct {
	purchaseRepo PurchaseRepository
}

func NewCheckout(repo PurchaseRepository) *Checkout {
	return &Checkout{repo}
}

func (c *Checkout) Purchase(cartItems []CartItem) (Invoice, error) {
	invoiceItems, totalCost, err := c.applyPromotions(cartItems)
	if err != nil {
		return Invoice{}, err
	}

	var reserveItems []CartItem
	for _, invItem := range invoiceItems {
		reserveItems = append(reserveItems, CartItem{
			SKU:      invItem.SKU,
			Quantity: invItem.Quantity,
		})
	}

	if err := c.purchaseRepo.BuyAllOrFail(reserveItems); err != nil {
		return Invoice{}, err
	}

	return Invoice{
		Items:     invoiceItems,
		TotalCost: totalCost,
	}, nil
}

func (c *Checkout) applyPromotions(items []CartItem) ([]InvoiceItem, float64, error) {
	var invoiceItems []InvoiceItem
	var totalCost float64

	for _, item := range items {
		product, err := c.purchaseRepo.GetProduct(item.SKU)
		if err != nil {
			return nil, 0, err
		}

		switch item.SKU {
		case "43N23P":
			macInvoice := InvoiceItem{
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
				freePi := InvoiceItem{
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
			chargedQty := (item.Quantity / 3) * 2 + (item.Quantity % 3)
			invoiceItem := InvoiceItem{
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
			invoiceItem := InvoiceItem{
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
			invoiceItem := InvoiceItem{
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
