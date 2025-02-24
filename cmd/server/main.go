package main

import (
	"log"
	"net/http"

	"checkout/internal/application"
	"checkout/internal/infrastructure/repository"
	"checkout/internal/interface/http/handlers"
)

func main() {
	repo := repository.NewDefaultMemoryDB()
	checkout := application.NewCheckout(repo)

	handler := handlers.NewPurchaseHandler(checkout)
	http.Handle("/purchase", handler)

	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server Failed: %v", err)
	}
}
