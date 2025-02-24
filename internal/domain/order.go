package domain

type OrderItem struct {
	SKU        string
	Name       string
	Quantity   int
	UnitPrice  float64
	TotalPrice float64
}

type Order struct {
	Items     []OrderItem
	TotalCost float64
}
