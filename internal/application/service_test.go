package application

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/assert"

	"checkout/internal/domain"
	"checkout/internal/infrastructure/repository"
)

func TestPurchase(t *testing.T) {
	repo := repository.NewEmptyMemoryDB()
	checkout := NewCheckout(repo)
	t.Run("testPurchase", func(t *testing.T) {
		testPurchase(t, repo, checkout)
	})

	repo = repository.NewEmptyMemoryDB()
	checkout = NewCheckout(repo)
	t.Run("testParrallelPurchase", func(t *testing.T) {
		testPartialPurchase(t, repo, checkout)
	})

}

func testPurchase(t *testing.T, repo *repository.MemoryDB, checkout *Checkout) {
	products := []*domain.Product{
		&domain.Product{SKU: "120P90", Name: "Google TV", Price: 49.99, Quantity: 10},
		&domain.Product{SKU: "A304SD", Name: "Alexa Speaker", Price: 109.50, Quantity: 10},
	}

	for _, product := range products {
		err := product.Validate()
		require.NoError(t, err)
		err = repo.UpdateProduct(*product)
		require.NoError(t, err)
	}

	items := []domain.CartItem{
		domain.CartItem{SKU: products[0].SKU, Quantity: 1},
		domain.CartItem{SKU: products[1].SKU, Quantity: 2},
	}

	invoice, err := checkout.Purchase(items)
	require.NoError(t, err)

	var requiredTotalCost float64 = 0
	requiredTotalCost += products[0].Price * float64(items[0].Quantity)
	requiredTotalCost += products[1].Price * float64(items[1].Quantity)

	assert.Equal(t, invoice.TotalCost, requiredTotalCost)
}

func testPartialPurchase(t *testing.T, repo *repository.MemoryDB, checkout *Checkout) {
	products := []*domain.Product{
		&domain.Product{SKU: "A", Name: "Item A", Price: 20.0, Quantity: 5},
		&domain.Product{SKU: "B", Name: "Item B", Price: 20.0, Quantity: 1},
	}

	for _, product := range products {
		err := product.Validate()
		require.NoError(t, err)
		err = repo.UpdateProduct(*product)
		require.NoError(t, err)
	}

	items := []domain.CartItem{
		domain.CartItem{SKU: "A", Quantity: 5},
		domain.CartItem{SKU: "B", Quantity: 2},
	}

	_, err := checkout.Purchase(items)
	require.Error(t, err)

	productA, _ := repo.GetProduct("A")
	assert.Equal(t, 5, productA.Quantity)
}


func TestPromotionMacBookPro(t *testing.T) {
	macBook := domain.Product{SKU: "43N23P", Name: "MacBook Pro", Price: 5399.99, Quantity: 5}
	pi := domain.Product{SKU: "234234", Name: "Raspberry Pi B", Price: 30.00, Quantity: 10}

	repo := repository.NewEmptyMemoryDB()
	require.NoError(t, repo.UpdateProduct(macBook))
	require.NoError(t, repo.UpdateProduct(pi))
	checkout := NewCheckout(repo)

	cart := []domain.CartItem{
		{SKU: "43N23P", Quantity: 1},
	}

	invoice, err := checkout.Purchase(cart)
	require.NoError(t, err)

	assert.Equal(t, 5399.99, invoice.TotalCost)

	// Check that both items appear in the invoice.
	var macItem, piItem *domain.InvoiceItem
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
