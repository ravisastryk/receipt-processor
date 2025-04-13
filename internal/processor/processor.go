package processor

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/fetch-rewards/receipt-processor/internal/models"
)

var (
	dateRegex = regexp.MustCompile(`^(\d{4})-(\d{2})-(\d{2})$`)
	timeRegex = regexp.MustCompile(`^(\d{2}):(\d{2})$`)
)

// CalculatePoints calculates the points for a receipt based on the rules
func CalculatePoints(receipt models.Receipt) int {
	points := 0

	// Rule 1: One point for every alphanumeric character in the retailer name
	retailerName := receipt.Retailer
	alphanumericCount := 0
	for _, char := range retailerName {
		if (char >= 'A' && char <= 'Z') ||
			(char >= 'a' && char <= 'z') ||
			(char >= '0' && char <= '9') {
			alphanumericCount++
		}
	}
	points += alphanumericCount

	// Rule 2: 50 points if the total is a round dollar amount with no cents
	total, _ := strconv.ParseFloat(receipt.Total, 64)
	if total == math.Floor(total) {
		points += 50
	}

	// Rule 3: 25 points if the total is a multiple of 0.25
	if math.Mod(total*100, 25) == 0 {
		points += 25
	}

	// Rule 4: 5 points for every two items on the receipt
	itemCount := len(receipt.Items)
	points += (itemCount / 2) * 5

	// Rule 5: If the trimmed length of the item description is a multiple of 3,
	// multiply the price by 0.2 and round up to the nearest integer
	for _, item := range receipt.Items {
		trimmedDesc := strings.TrimSpace(item.ShortDescription)
		if len(trimmedDesc)%3 == 0 {
			price, _ := strconv.ParseFloat(item.Price, 64)
			points += int(math.Ceil(price * 0.2))
		}
	}

	// Rule 6: 6 points if the purchase date is odd
	matches := dateRegex.FindStringSubmatch(receipt.PurchaseDate)
	if len(matches) == 4 {
		day, err := strconv.Atoi(matches[3])
		if err != nil {
			// TODO: Errors can be captured and handled at the caller
			fmt.Println("eror converting string to int at matches[3]", err)
			return -1
		}
		if day%2 == 1 {
			points += 6
		}
	}

	// Rule 7: 10 points if the purchase time is after 2:00pm and before 4:00pm
	timeMatches := timeRegex.FindStringSubmatch(receipt.PurchaseTime)
	if len(timeMatches) == 3 {
		hour, err := strconv.Atoi(timeMatches[1])
		if err != nil {
			// TODO: Errors can be captured and handled at the caller
			fmt.Println("eror converting string to int at timeMatches[1]", err)
			return -1
		}

		minute, err := strconv.Atoi(timeMatches[2])
		if err != nil {
			// TODO: Errors can be captured and handled at the caller
			fmt.Println("eror converting string to int at timeMatches[2]", err)
			return -1
		}

		purchaseTime := 60*hour + minute
		if purchaseTime > 14*60 && purchaseTime < 16*60 {
			points += 10
		}
	}

	return points
}
