package application

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/assert"
)

func TestPurchase(t *testing.T) {
	repository := NewEmptyMemoryDB()
	checkout := NewCheckout(repository)
	t.Run("testPurchase", func(t *testing.T) {
		testPurchase(t, repository, checkout)
	})

	repository = NewEmptyMemoryDB()
	checkout = NewCheckout(repository)
	t.Run("testParrallelPurchase", func(t *testing.T) {
		testPartialPurchase(t, repository, checkout)
	})

}

func testPurchase(t *testing.T, repository *MemoryDB, checkout *Checkout) {
	products := []*Product{
		&Product{SKU: "120P90", Name: "Google TV", Price: 49.99, Quantity: 10},
		&Product{SKU: "A304SD", Name: "Alexa Speaker", Price: 109.50, Quantity: 10},
	}

	for _, product := range products {
		err := product.Validate()
		require.NoError(t, err)
		err = repository.UpdateProduct(*product)
		require.NoError(t, err)
	}

	items := []CartItem{
		CartItem{SKU: products[0].SKU, Quantity: 1},
		CartItem{SKU: products[1].SKU, Quantity: 2},
	}

	invoice, err := checkout.Purchase(items)
	require.NoError(t, err)

	var requiredTotalCost float64 = 0
	requiredTotalCost += products[0].Price * float64(items[0].Quantity)
	requiredTotalCost += products[1].Price * float64(items[1].Quantity)

	assert.Equal(t, invoice.TotalCost, requiredTotalCost)
}

func testPartialPurchase(t *testing.T, repository *MemoryDB, checkout *Checkout) {
	products := []*Product{
		&Product{SKU: "A", Name: "Item A", Price: 20.0, Quantity: 5},
		&Product{SKU: "B", Name: "Item B", Price: 20.0, Quantity: 1},
	}

	for _, product := range products {
		err := product.Validate()
		require.NoError(t, err)
		err = repository.UpdateProduct(*product)
		require.NoError(t, err)
	}

	items := []CartItem{
		CartItem{SKU: "A", Quantity: 5},
		CartItem{SKU: "B", Quantity: 2},
	}

	_, err := checkout.Purchase(items)
	require.Error(t, err)

	productA, _ := repository.GetProduct("A")
	assert.Equal(t, 5, productA.Quantity)

}


func TestPromotionMacBookPro(t *testing.T) {
	macBook := Product{SKU: "43N23P", Name: "MacBook Pro", Price: 5399.99, Quantity: 5}
	pi := Product{SKU: "234234", Name: "Raspberry Pi B", Price: 30.00, Quantity: 10}

	repo := NewEmptyMemoryDB()
	require.NoError(t, repo.UpdateProduct(macBook))
	require.NoError(t, repo.UpdateProduct(pi))
	checkout := NewCheckout(repo)

	cart := []CartItem{
		{SKU: "43N23P", Quantity: 1},
	}

	invoice, err := checkout.Purchase(cart)
	require.NoError(t, err)

	assert.Equal(t, 5399.99, invoice.TotalCost)

	// Check that both items appear in the invoice.
	var macItem, piItem *InvoiceItem
	for i, item := range invoice.Items {
		if item.SKU == "43N23P" {
			macItem = &invoice.Items[i]
		} else if item.SKU == "234234" {
			piItem = &invoice.Items[i]
		}
	}

	require.NotNil(t, macItem, "Invoice should include MacBook Pro")
	require.NotNil(t, piItem, "Invoice should include free Raspberry Pi B")

	assert.Equal(t, 1, macItem.Quantity)
	assert.Equal(t, 5399.99, macItem.TotalPrice)
	assert.Equal(t, 1, piItem.Quantity)
	assert.Equal(t, 0.0, piItem.TotalPrice)
}
