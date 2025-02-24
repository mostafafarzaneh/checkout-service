package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"checkout/internal/application"
	"checkout/internal/domain"
)

type PurchaseRequest struct {
	Items []domain.CartItem `json:"items"`
}

type PurchaseResponse struct {
	Invoice domain.Invoice `json:"invoice"`
	Error   string         `json:"error, omitempty"`
}

type PurchaseHandler struct {
	Checkout *application.Checkout
}

func NewPurchaseHandler(checkout *application.Checkout) *PurchaseHandler {
	return &PurchaseHandler{checkout}
}

func (h *PurchaseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req PurchaseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	invoice, err := h.Checkout.Purchase(req.Items)
	if err != nil {
		log.Printf("Purchase error: %v", err)
		resp := PurchaseResponse{
			Error:   err.Error(),
			Invoice: domain.Invoice{},
		}
		w.WriteHeader(http.StatusConflict)
		_ = json.NewEncoder(w).Encode(resp)
		return
	}

	resp := PurchaseResponse{
		Invoice: invoice,
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}
