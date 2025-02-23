package application

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/assert"
)

func TestPurchase(t *testing.T) {
	repository := NewEmptyMemoryDB()
	t.Run("testUpdateProduct", func(t *testing.T) {
		testUpdateProduct(t, repository)
	})

	t.Run("testBuyProduct", func(t *testing.T) {
		testBuyProduct(t, repository)
	})
}

func testUpdateProduct(t *testing.T, repository *MemoryDB) {
	product := &Product{SKU: "120P90", Name: "Google TV", Price: 49.99, Quantity: 10}
	err := product.Validate()
	require.NoError(t, err)
	err = repository.UpdateProduct(product)
	require.NoError(t, err)

	gotProduct, err := repository.GetProduct(product.SKU)
	require.NoError(t, err)

	assert.Equal(t, product, gotProduct)
}

func testBuyProduct(t *testing.T, repository *MemoryDB) {
	product := &Product{SKU: "120P90", Name: "Google TV", Price: 49.99, Quantity: 10}
	err := product.Validate()
	require.NoError(t, err)
	err = repository.UpdateProduct(product)
	require.NoError(t, err)

	err = repository.BuyProduct(product.SKU, 5)
	require.NoError(t, err)

	requiredProduct := &Product{SKU: "120P90", Name: "Google TV", Price: 49.99, Quantity: 5}
	gotProduct, err := repository.GetProduct(product.SKU)
	require.NoError(t, err)
	assert.Equal(t, requiredProduct, gotProduct)
}
