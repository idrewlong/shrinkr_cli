APP_NAME := shrinkr
VERSION := 1.0.0

.PHONY: build install clean run

# Build the binary
build:
	go build -ldflags "-s -w" -o $(APP_NAME) .

# Install globally
install:
	go install .

# Clean build artifacts
clean:
	rm -f $(APP_NAME)

# Run with sample images
run:
	go run . ../images

# Build for current platform with version info
release:
	go build -ldflags "-s -w -X main.Version=$(VERSION)" -o $(APP_NAME) .
