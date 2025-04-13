.PHONY: build run test docker-build docker-run clean

all: build

build:
	go build -o bin/receipt-processor ./cmd/api

run:
	go run ./cmd/api

test:
	go test ./...


docker-build:
	docker build -t receipt-processor .

docker-run:
	docker run -p 8080:8080 receipt-processor

clean:
	rm -rf bin/
