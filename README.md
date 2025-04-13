# Receipt Processor Challenge

This repository contains a solution to the [Fetch Rewards Receipt Processor Challenge](https://github.com/fetch-rewards/receipt-processor-challenge).

## Project Structure

```
$tree
receipt-processor/
.
├── Dockerfile # Container definition
├── Makefile
├── README.md
├── api
│   └── openapi.yaml
├── cmd
│   └── api
│       └── main.go # Application entry point
├── go.mod
├── go.sum
└── internal
    ├── handlers  # HTTP request handlers
    │   ├── handlers.go
    │   └── handlers_test.go
    ├── models # Data models
    │   └── models.go
    ├── processor  # Receipt point processing logic
    │   ├── processor.go
    │   └── processor_test.go
    └── storage # Data storage interfaces and implementations
        ├── memory.go
        └── storage.go

9 directories, 14 files
```

### Setup

1. Make sure you have Go 1.21+ installed
2. Clone the repository
3. Install dependencies:
   ```bash
   go mod download
   ```
4. Run the application:
   ```bash
   go run cmd/api/main.go
   ```

## API Endpoints

The API provides two main endpoints:

### Process a Receipt

```
POST /receipts/process
```

Request body example:
```json
{
  "retailer": "Target",
  "purchaseDate": "2022-01-01",
  "purchaseTime": "13:01",
  "items": [
    {
      "shortDescription": "Mountain Dew 12PK",
      "price": "6.49"
    },
    {
      "shortDescription": "Emils Cheese Pizza",
      "price": "12.25"
    }
  ],
  "total": "18.74"
}
```

Response:
```json
{
  "id": "7fb1377b-b223-49d9-a31a-5a02701dd310"
}
```

### Get Points for a Receipt

```
GET /receipts/{id}/points
```

Response:
```json
{
  "points": 28
}
```

## Point Calculation Rules

Points are calculated based on these rules:

1. One point for every alphanumeric character in the retailer name
2. 50 points if the total is a round dollar amount with no cents
3. 25 points if the total is a multiple of 0.25
4. 5 points for every two items on the receipt
5. If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and round up to the nearest integer
6. 6 points if the purchase date is odd
7. 10 points if the purchase time is after 2:00pm and before 4:00pm

## Development

### Running Tests

Run all tests:
```bash
go test ./...
```

Run with coverage:
```bash
go test ./... -cover
```

### Building and Running with Docker

Build the Docker image:
```bash
docker build -t receipt-processor .
```

Run the container:
```bash
docker run -p 8080:8080 receipt-processor
```
### Sample Input & Output

## Service running in a local terminal

<img width="646" alt="Screenshot 2025-04-13 at 4 45 06 PM" src="https://github.com/user-attachments/assets/a0275114-75b1-4302-9985-909987935b73" />

## API Outputs

<img width="723" alt="Screenshot 2025-04-13 at 4 44 23 PM" src="https://github.com/user-attachments/assets/ac270666-4c3f-40cc-9ecd-46ab5dba6b80" />

