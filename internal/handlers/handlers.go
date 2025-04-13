package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/fetch-rewards/receipt-processor/internal/models"
	"github.com/fetch-rewards/receipt-processor/internal/processor"
	"github.com/fetch-rewards/receipt-processor/internal/storage"
)

// Handler contains the HTTP handlers and dependencies
type Handler struct {
	Storage storage.ReceiptStorage
}

// NewHandler creates a new Handler instance
func NewHandler(storage storage.ReceiptStorage) *Handler {
	return &Handler{
		Storage: storage,
	}
}

// ProcessReceipt handles the POST request to process a receipt
func (h *Handler) ProcessReceipt(w http.ResponseWriter, r *http.Request) {
	var receipt models.Receipt

	// Decode JSON request
	err := json.NewDecoder(r.Body).Decode(&receipt)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Generate a unique ID
	id := uuid.New().String()

	// Calculate points
	points := processor.CalculatePoints(receipt)

	// Store receipt ID and points
	if err := h.Storage.SaveReceipt(id, points); err != nil {
		http.Error(w, "Failed to process receipt", http.StatusInternalServerError)
		return
	}

	// Return response
	response := models.ReceiptResponse{ID: id}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetPoints handles the GET request to retrieve points for a receipt
func (h *Handler) GetPoints(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Check if receipt exists
	points, exists := h.Storage.GetPoints(id)
	if !exists {
		http.Error(w, "Receipt not found", http.StatusNotFound)
		return
	}

	// Return response
	response := models.PointsResponse{Points: points}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
