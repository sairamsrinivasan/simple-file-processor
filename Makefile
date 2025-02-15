BINARY_NAME = simple-file-processor

# Default target to build the application
all: build

build: 
	@echo "Building the application..."
	go build -o ./bin/$(BINARY_NAME) main.go
	@echo "Build complete."

run:
	@echo "Running the application..."
	./bin/$(BINARY_NAME)

clean:
	@echo "Cleaning up..."
	rm -rf ./bin/$(BINARY_NAME)
	@echo "Clean complete."

test:
	@echo "Running tests..."
	go test ./...
	@echo "Tests complete."	