APP_NAME := shrinkr
VERSION := 1.0.0

.PHONY: build install clean run test bench release-local tag-release

# Build the binary
build:
	go build -ldflags "-s -w" -o $(APP_NAME) .

# Install globally
install:
	go install .

# Clean build artifacts
clean:
	rm -f $(APP_NAME)
	rm -rf dist/

# Run with sample images
run:
	go run . ../images

# Run all tests
test:
	go test ./... -v -timeout 120s

# Run format benchmarks (shows per-format encoding speed)
bench:
	go test ./internal/compressor/ -bench=. -benchtime=3x -timeout 120s

# Test GoReleaser config locally (no publish)
release-local:
	goreleaser release --snapshot --clean

# Tag and push a release (triggers GitHub Actions)
# Usage: make tag-release V=1.0.0
tag-release:
	@if [ -z "$(V)" ]; then echo "Usage: make tag-release V=1.0.0"; exit 1; fi
	git tag -a v$(V) -m "Release v$(V)"
	git push origin v$(V)
