package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/fetch-rewards/receipt-processor/internal/handlers"
	"github.com/fetch-rewards/receipt-processor/internal/storage"
)

func main() {
	// Create a new router
	r := mux.NewRouter()

	// Initialize storage
	store := storage.NewMemoryStorage()

	// Initialize handlers
	handler := handlers.NewHandler(store)

	// Define routes
	// TODO Authentication middleware is avoided for simplicity
	r.HandleFunc("/receipts/process", handler.ProcessReceipt).Methods("POST")
	r.HandleFunc("/receipts/{id}/points", handler.GetPoints).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	fmt.Printf("Server started on :%s\n", port)
	// TOO graceful shutdown logic can be done as an enhancement
	log.Fatal(http.ListenAndServe(":"+port, r))
}
