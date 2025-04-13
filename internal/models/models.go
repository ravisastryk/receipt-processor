package models

// Receipt represents the receipt data structure
type Receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []Item `json:"items"`
	Total        string `json:"total"`
}

// Item represents a single item in the receipt
type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

// ReceiptResponse is the response for a processed receipt
type ReceiptResponse struct {
	ID string `json:"id"`
}

// PointsResponse is the response for a points inquiry
type PointsResponse struct {
	Points int `json:"points"`
}
