package processor

import (
	"testing"

	"github.com/fetch-rewards/receipt-processor/internal/models"
)

func TestCalculatePoints(t *testing.T) {
	tests := []struct {
		name    string
		receipt models.Receipt
		want    int
	}{
		{
			name: "Example 1",
			receipt: models.Receipt{
				Retailer:     "Target",
				PurchaseDate: "2022-01-01",
				PurchaseTime: "13:01",
				Items: []models.Item{
					{ShortDescription: "Mountain Dew 12PK", Price: "6.49"},
					{ShortDescription: "Emils Cheese Pizza", Price: "12.25"},
					{ShortDescription: "Knorr Creamy Chicken", Price: "1.26"},
					{ShortDescription: "Doritos Nacho Cheese", Price: "3.35"},
					{ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ", Price: "12.00"},
				},
				Total: "35.35",
			},
			want: 28,
		},
		{
			name: "Example 2",
			receipt: models.Receipt{
				Retailer:     "M&M Corner Market",
				PurchaseDate: "2022-03-20",
				PurchaseTime: "14:33",
				Items: []models.Item{
					{ShortDescription: "Gatorade", Price: "2.25"},
					{ShortDescription: "Gatorade", Price: "2.25"},
					{ShortDescription: "Gatorade", Price: "2.25"},
					{ShortDescription: "Gatorade", Price: "2.25"},
				},
				Total: "9.00",
			},
			want: 109,
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CalculatePoints(tt.receipt); got != tt.want {
				t.Errorf("CalculatePoints() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test the individual point calculation rules
func TestPointCalculationRules(t *testing.T) {
	// Test Rule 1: One point for every alphanumeric character in the retailer name
	retailerTest := models.Receipt{Retailer: "Target"}
	if points := CalculatePoints(retailerTest); points != 6 {
		t.Errorf("Rule 1 failed: got %v, want 6", points)
	}

	// Test Rule 2: 50 points if the total is a round dollar amount with no cents
	roundDollarTest := models.Receipt{Retailer: "", Total: "100.00"}
	if points := CalculatePoints(roundDollarTest); points != 50 {
		t.Errorf("Rule 2 failed: got %v, want 50", points)
	}

	// Test Rule 3: 25 points if the total is a multiple of 0.25
	quarterMultipleTest := models.Receipt{Retailer: "", Total: "10.75"}
	if points := CalculatePoints(quarterMultipleTest); points != 25 {
		t.Errorf("Rule 3 failed: got %v, want 25", points)
	}

	// Test Rule 4: 5 points for every two items on the receipt
	itemCountTest := models.Receipt{
		Retailer: "",
		Items: []models.Item{
			{ShortDescription: "Item 1", Price: "1.00"},
			{ShortDescription: "Item 2", Price: "2.00"},
			{ShortDescription: "Item 3", Price: "3.00"},
			{ShortDescription: "Item 4", Price: "4.00"},
		},
	}
	expectedItemPoints := 10 // 4 items = 2 pairs = 10 points
	if points := CalculatePoints(itemCountTest); points != expectedItemPoints {
		t.Errorf("Rule 4 failed: got %v, want %v", points, expectedItemPoints)
	}

	// Test Rule 5: Trimmed length multiple of 3
	trimmedLengthTest := models.Receipt{
		Retailer: "",
		Items: []models.Item{
			{ShortDescription: "ABC", Price: "10.00"}, // 3 chars = multiple of 3, 10.00 * 0.2 = 2 points
		},
	}
	expectedTrimPoints := 2
	if points := CalculatePoints(trimmedLengthTest); points != expectedTrimPoints {
		t.Errorf("Rule 5 failed: got %v, want %v", points, expectedTrimPoints)
	}

	// Test Rule 6: 6 points if the purchase date is odd
	oddDateTest := models.Receipt{Retailer: "", PurchaseDate: "2022-01-01"}
	if points := CalculatePoints(oddDateTest); points != 6 {
		t.Errorf("Rule 6 failed: got %v, want 6", points)
	}

	// Test Rule 7: 10 points if the purchase time is after 2:00pm and before 4:00pm
	timeTest := models.Receipt{Retailer: "", PurchaseTime: "14:30"}
	if points := CalculatePoints(timeTest); points != 10 {
		t.Errorf("Rule 7 failed: got %v, want 10", points)
	}
}
