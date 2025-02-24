package domain

type InvoiceItem struct {
	SKU        string
	Name       string
	Quantity   int
	UnitPrice  float64
	TotalPrice float64
}

type Invoice struct {
	Items     []InvoiceItem
	TotalCost float64
}
