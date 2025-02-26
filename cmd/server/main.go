package main

import (
	"log"
	"net/http"

	"checkout/internal/application"
	"checkout/internal/infrastructure/payment"
	"checkout/internal/infrastructure/orderprocessing"
	"checkout/internal/infrastructure/repository"
	"checkout/internal/interface/http/handlers"
)

func main() {
	paymentProcessor := &payment.DummyPaymentProcessor{}
	orderProcessor := &orderprocessing.DummyOrderProcessor{}

	repo := repository.NewDefaultMemoryDB()
	checkout := application.NewCheckout(repo, paymentProcessor, orderProcessor)

	handler := handlers.NewPurchaseHandler(checkout)
	http.Handle("/purchase", handler)

	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server Failed: %v", err)
	}
}
