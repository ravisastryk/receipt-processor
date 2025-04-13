package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/fetch-rewards/receipt-processor/internal/models"
	"github.com/fetch-rewards/receipt-processor/internal/storage"
)

func TestProcessReceipt(t *testing.T) {
	// Setup
	store := storage.NewMemoryStorage()
	handler := NewHandler(store)

	// Create test receipt
	receipt := models.Receipt{
		Retailer:     "Target",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "13:01",
		Items: []models.Item{
			{ShortDescription: "Mountain Dew 12PK", Price: "6.49"},
			{ShortDescription: "Emils Cheese Pizza", Price: "12.25"},
		},
		Total: "18.74",
	}

	// Convert to JSON
	receiptJSON, err := json.Marshal(receipt)
	if err != nil {
		t.Fatalf("Failed to marshal receipt: %v", err)
	}

	// Create request
	req, err := http.NewRequest("POST", "/receipts/process", bytes.NewBuffer(receiptJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Execute request
	rr := httptest.NewRecorder()
	http.HandlerFunc(handler.ProcessReceipt).ServeHTTP(rr, req)

	// Check status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check response
	var response models.ReceiptResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.ID == "" {
		t.Errorf("Expected non-empty ID")
	}
}

func TestGetPoints(t *testing.T) {
	// Setup
	store := storage.NewMemoryStorage()
	handler := NewHandler(store)

	// Add a receipt to the store
	testID := "test-id"
	testPoints := 100
	store.SaveReceipt(testID, testPoints)

	// Create request
	req, err := http.NewRequest("GET", "/receipts/"+testID+"/points", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Set up the router to capture variables
	router := mux.NewRouter()
	router.HandleFunc("/receipts/{id}/points", handler.GetPoints)

	// Execute request
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Check status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check response
	var response models.PointsResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Points != testPoints {
		t.Errorf("Expected points %v, got %v", testPoints, response.Points)
	}
}
