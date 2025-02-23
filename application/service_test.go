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
}

func testPurchase(t *testing.T, repository *MemoryDB, checkout *Checkout) {
	products := []*Product{
		&Product{SKU: "120P90", Name: "Google TV", Price: 49.99, Quantity: 10},
		&Product{SKU: "43N23P", Name: "MacBook Pro", Price: 5399.99, Quantity: 5},
		&Product{SKU: "A304SD", Name: "Alexa Speaker", Price: 109.50, Quantity: 10},
	}

	for _, product := range products {
		err := product.Validate()
		require.NoError(t, err)
		err = repository.UpdateProduct(product)
		require.NoError(t, err)
	}

	items := []CartItem{
		CartItem{SKU: products[0].SKU, Quantity: 5},
		CartItem{SKU: products[1].SKU, Quantity: 2},
	}

	totalCost, err := checkout.Purchase(items)
	require.NoError(t, err)

	var requiredTotalCost float64 = 0
	requiredTotalCost += products[0].Price * float64(items[0].Quantity)
	requiredTotalCost += products[1].Price * float64(items[1].Quantity)

	assert.Equal(t, totalCost, requiredTotalCost)
}
