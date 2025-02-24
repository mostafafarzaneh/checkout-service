package repository

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"

	"checkout/internal/application/port"
	"checkout/internal/domain"
)

func TestMemoryDB(t *testing.T) {
	repository := NewEmptyMemoryDB()
	t.Run("testUpdateProduct", func(t *testing.T) {
		testUpdateProduct(t, repository)
	})

	t.Run("testBuyProduct", func(t *testing.T) {
		testBuyProduct(t, repository)
	})
}

func testUpdateProduct(t *testing.T, repository port.PurchaseRepository) {
	product := &domain.Product{SKU: "120P90", Name: "Google TV", Price: 49.99, Quantity: 10}
	err := product.Validate()
	require.NoError(t, err)
	err = repository.UpdateProduct(*product)
	require.NoError(t, err)

	gotProduct, err := repository.GetProduct(product.SKU)
	require.NoError(t, err)

	assert.Equal(t, *product, gotProduct)
}

func testBuyProduct(t *testing.T, repository port.PurchaseRepository) {
	product := &domain.Product{SKU: "120P90", Name: "Google TV", Price: 49.99, Quantity: 10}
	err := product.Validate()
	require.NoError(t, err)
	err = repository.UpdateProduct(*product)
	require.NoError(t, err)

	err = repository.BuyAllOrFail([]domain.CartItem{domain.CartItem{product.SKU, 5}})
	require.NoError(t, err)

	requiredProduct := &domain.Product{SKU: "120P90", Name: "Google TV", Price: 49.99, Quantity: 5}
	gotProduct, err := repository.GetProduct(product.SKU)
	require.NoError(t, err)
	assert.Equal(t, *requiredProduct, gotProduct)
}

func TestConcurrentPurchase(t *testing.T) {
	repo := NewEmptyMemoryDB()
	_ = repo.UpdateProduct(domain.Product{
		SKU:      "120P90",
		Name:     "Google TV",
		Price:    49.99,
		Quantity: 100,
	})

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := repo.BuyAllOrFail([]domain.CartItem{domain.CartItem{"120P90", 5}})
			require.NoError(t, err)
		}()
	}
	wg.Wait()

	product, err := repo.GetProduct("120P90")
	require.NoError(t, err)
	assert.Equal(t, 50, product.Quantity)
}
