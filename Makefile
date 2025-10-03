.PHONY: run build tidy

run:
	@echo "Running the application..."
	go run ./cmd/server

build:
	@echo "Building the application..."
	go build -o app cmd/server/*.go

tidy:
	@echo "Tidying dependencies..."
	go mod tidy
