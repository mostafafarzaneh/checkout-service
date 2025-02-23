package application

import (
	"errors"
	"sync"
)

type InventoryItem struct {
	SKU string
	Name string
	Price float64
	Quantity int
}

type MemoryDB struct {
	inventory map[string]InventoryItem
	mu sync.Mutex
}

var (
	NotEnoughQuantityErr = errors.New("Not enough quantity")
	InventoryNotFoundErr = errors.New("Inventory Not found")
)

func NewEmptyMemoryDB() *MemoryDB {
	inventory := make(map[string]InventoryItem)
	return &MemoryDB{inventory: inventory}
}

func NewDefaultMemoryDB() *MemoryDB {
	inventory := make(map[string]InventoryItem)
	inventory["120P90"] = InventoryItem{SKU: "120P90", Name: "Google TV", Price: 49.99, Quantity: 10}
	inventory["43N23P"] = InventoryItem{SKU: "43N23P", Name: "MacBook Pro", Price: 5399.99, Quantity: 5}
	inventory["A304SD"] = InventoryItem{SKU: "A304SD", Name: "Alexa Speaker", Price: 109.50, Quantity: 10}
	inventory["234234"] = InventoryItem{SKU: "234234", Name: "Raspberry Pi B", Price: 30.00, Quantity: 2}
	return &MemoryDB{inventory: inventory}
}

func (m *MemoryDB) GetProduct(sku string) (*Product, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	inventory, ok := m.inventory[sku]
	if !ok {
		return nil, InventoryNotFoundErr
	}
	return convertInventoryItemToProduct(inventory), nil
}

func (m *MemoryDB) UpdateProduct(product *Product) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	sku := product.SKU
	name := product.Name
	price := product.Price
	quantity := product.Quantity
	m.inventory[sku] = InventoryItem{SKU: sku, Name: name, Price: price, Quantity: quantity}
	return nil
}

func (m *MemoryDB) BuyProduct(sku string, quantity int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.inventory[sku].Quantity < quantity {
		return NotEnoughQuantityErr
	}
	newItem := m.inventory[sku]
	newItem.Quantity -= quantity
	m.inventory[sku] = newItem
	return nil
}

func convertInventoryItemToProduct(item InventoryItem) *Product {
	return &Product{
		SKU: item.SKU,
		Name: item.Name,
		Price: item.Price,
		Quantity: item.Quantity,
	}
}
