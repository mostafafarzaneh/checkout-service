package repository

import (
	"errors"
	"sync"

	"checkout/internal/domain"
)

type InventoryItem struct {
	SKU      string
	Name     string
	Price    float64
	Quantity int
}

type MemoryDB struct {
	inventory map[string]InventoryItem
	mu        sync.Mutex
}

var (
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

func (m *MemoryDB) GetProduct(sku string) (domain.Product, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	inventory, ok := m.inventory[sku]
	if !ok {
		return domain.Product{}, InventoryNotFoundErr
	}
	return convertInventoryItemToProduct(inventory), nil
}

func (m *MemoryDB) UpdateProduct(product domain.Product) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	sku := product.SKU
	name := product.Name
	price := product.Price
	quantity := product.Quantity
	m.inventory[sku] = InventoryItem{SKU: sku, Name: name, Price: price, Quantity: quantity}
	return nil
}

func (m *MemoryDB) BuyAllOrFail(items []domain.CartItem) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	productList := []domain.Product{}
	for _, item := range items {
		sku := item.SKU
		p := convertInventoryItemToProduct(m.inventory[sku])
		if err := p.BuyProduct(item); err != nil {
			return err
		}
		productList = append(productList, p)
	}
	for _, product := range productList {
		m.inventory[product.SKU] = convertProductToInventoryItem(product)
	}
	return nil
}

func (m *MemoryDB) RestoreProducts(items []domain.CartItem) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, item := range items {
		invItem, ok := m.inventory[item.SKU]
		if !ok {
			continue
		}
		invItem.Quantity += item.Quantity
		m.inventory[item.SKU] = invItem
	}
	return nil
}

func convertInventoryItemToProduct(item InventoryItem) domain.Product {
	return domain.Product{
		SKU:      item.SKU,
		Name:     item.Name,
		Price:    item.Price,
		Quantity: item.Quantity,
	}
}

func convertProductToInventoryItem(product domain.Product) InventoryItem {
	return InventoryItem{
		SKU:      product.SKU,
		Name:     product.Name,
		Price:    product.Price,
		Quantity: product.Quantity,
	}
}
