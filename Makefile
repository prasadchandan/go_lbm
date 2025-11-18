.PHONY: all build clean test install-deps install-gomobile android ios help

# Default target
all: build

# Install development dependencies
install-deps:
	@echo "Installing development dependencies..."
	go mod download
	go mod verify

# Install gomobile for mobile builds
install-gomobile:
	@echo "Installing gomobile..."
	go install golang.org/x/mobile/cmd/gomobile@latest
	go install golang.org/x/mobile/cmd/gobind@latest
	gomobile init

# Build for current platform
build:
	@echo "Building for current platform..."
	go build -v -o go_lbm .

# Build for Linux x86_64
build-linux-amd64:
	@echo "Building for Linux x86_64..."
	GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -v -o go_lbm-linux-amd64 .

# Build for Linux ARM64
build-linux-arm64:
	@echo "Building for Linux ARM64..."
	GOOS=linux GOARCH=arm64 CGO_ENABLED=1 CC=aarch64-linux-gnu-gcc go build -v -o go_lbm-linux-arm64 .

# Build for Windows x86_64
build-windows-amd64:
	@echo "Building for Windows x86_64..."
	GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build -v -o go_lbm-windows-amd64.exe .

# Build for macOS x86_64
build-darwin-amd64:
	@echo "Building for macOS x86_64..."
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -v -o go_lbm-darwin-amd64 .

# Build for macOS ARM64
build-darwin-arm64:
	@echo "Building for macOS ARM64..."
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 go build -v -o go_lbm-darwin-arm64 .

# Build for Android
android: install-gomobile
	@echo "Building for Android..."
	gomobile build -target=android -o go_lbm.apk .

# Build for iOS
ios: install-gomobile
	@echo "Building for iOS..."
	gomobile build -target=ios -o go_lbm.app .

# Run the application
run: build
	./go_lbm

# Run tests
test:
	go test -v -race ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f go_lbm go_lbm-* *.apk *.app
	rm -rf *.dSYM

# Show help
help:
	@echo "Available targets:"
	@echo "  make build              - Build for current platform"
	@echo "  make build-linux-amd64  - Build for Linux x86_64"
	@echo "  make build-linux-arm64  - Build for Linux ARM64"
	@echo "  make build-windows-amd64- Build for Windows x86_64"
	@echo "  make build-darwin-amd64 - Build for macOS x86_64"
	@echo "  make build-darwin-arm64 - Build for macOS ARM64"
	@echo "  make android            - Build APK for Android"
	@echo "  make ios                - Build for iOS"
	@echo "  make run                - Build and run"
	@echo "  make test               - Run tests"
	@echo "  make clean              - Clean build artifacts"
	@echo "  make install-deps       - Install Go dependencies"
	@echo "  make install-gomobile   - Install gomobile tools"
