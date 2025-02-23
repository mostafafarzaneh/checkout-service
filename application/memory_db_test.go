package application

import (
	"testing"
	"sync"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/assert"
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

func testUpdateProduct(t *testing.T, repository *MemoryDB) {
	product := &Product{SKU: "120P90", Name: "Google TV", Price: 49.99, Quantity: 10}
	err := product.Validate()
	require.NoError(t, err)
	err = repository.UpdateProduct(*product)
	require.NoError(t, err)

	gotProduct, err := repository.GetProduct(product.SKU)
	require.NoError(t, err)

	assert.Equal(t, *product, gotProduct)
}

func testBuyProduct(t *testing.T, repository *MemoryDB) {
	product := &Product{SKU: "120P90", Name: "Google TV", Price: 49.99, Quantity: 10}
	err := product.Validate()
	require.NoError(t, err)
	err = repository.UpdateProduct(*product)
	require.NoError(t, err)

	totalCost, err := repository.BuyProducts([]CartItem{CartItem{product.SKU, 5}})
	require.NoError(t, err)
	require.Equal(t, 49.99*float64(5), totalCost)

	requiredProduct := &Product{SKU: "120P90", Name: "Google TV", Price: 49.99, Quantity: 5}
	gotProduct, err := repository.GetProduct(product.SKU)
	require.NoError(t, err)
	assert.Equal(t, *requiredProduct, gotProduct)
}

func TestConcurrentPurchase(t *testing.T) {
    repo := NewEmptyMemoryDB()
    _ = repo.UpdateProduct(Product{
        SKU: "120P90",
        Name: "Google TV",
        Price: 49.99,
        Quantity: 100,
    })

    var wg sync.WaitGroup
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            totalCost, err := repo.BuyProducts([]CartItem{CartItem{"120P90", 5}})
            require.NoError(t, err)
	    require.Equal(t, 49.99*float64(5), totalCost)
        }()
    }
    wg.Wait()

    product, err := repo.GetProduct("120P90")
    require.NoError(t, err)
    assert.Equal(t, 50, product.Quantity)
}
