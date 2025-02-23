BINARY_NAME = simple-file-processor

# Default target to build the application
all: 
	@echo "Starting the build process..."
	build
	@echo "Build process complete. To run the application, use 'make run'."
	@echo "To clean up, use 'make clean'."

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