## Lattice Boltzman Solver

A cross-platform Lattice Boltzmann fluid dynamics solver written in Go. It uses [gomobile](https://github.com/golang/mobile) to build artifacts for Android and iOS, and can also run on macOS, Linux, and Windows.

[![Build and Release](https://github.com/prasadchandan/go_lbm/actions/workflows/build.yml/badge.svg)](https://github.com/prasadchandan/go_lbm/actions/workflows/build.yml)
[![Continuous Integration](https://github.com/prasadchandan/go_lbm/actions/workflows/ci.yml/badge.svg)](https://github.com/prasadchandan/go_lbm/actions/workflows/ci.yml)

### Features

- Cross-platform fluid dynamics simulation using Lattice Boltzmann method
- Supports Windows, macOS, Linux (x86_64 and ARM64)
- Mobile support for Android and iOS
- Interactive UI with touch/mouse controls
- Real-time visualization with OpenGL ES

### Prerequisites

- **Go 1.24+** - [Install Go](https://golang.org/doc/install)
- **Platform-specific dependencies:**
  - **Linux**: `sudo apt-get install libegl1-mesa-dev libgles2-mesa-dev libx11-dev`
  - **macOS**: Xcode Command Line Tools
  - **Windows**: No additional dependencies for basic build
  - **Android**: Android SDK and NDK (automatically handled by gomobile)
  - **iOS**: Xcode (macOS only)

### Quick Start

#### Building for Your Current Platform

```bash
# Clone the repository
git clone https://github.com/prasadchandan/go_lbm.git
cd go_lbm

# Install dependencies
go mod download

# Build
go build .

# Run
./go_lbm
```

#### Using Makefile

The project includes a comprehensive Makefile for easy building:

```bash
# Build for current platform
make build

# Build for specific platforms
make build-linux-amd64    # Linux x86_64
make build-linux-arm64    # Linux ARM64
make build-windows-amd64  # Windows x86_64
make build-darwin-amd64   # macOS Intel
make build-darwin-arm64   # macOS Apple Silicon

# Build mobile apps
make android              # Build APK for Android
make ios                  # Build for iOS (macOS only)

# Other commands
make run                  # Build and run
make test                 # Run tests
make clean                # Clean build artifacts
make help                 # Show all available targets
```

### Building for Mobile

#### Android

```bash
# Install gomobile (one-time setup)
go install golang.org/x/mobile/cmd/gomobile@latest
go install golang.org/x/mobile/cmd/gobind@latest
gomobile init

# Build APK
gomobile build -target=android -o go_lbm.apk .
```

#### iOS

```bash
# Install gomobile (one-time setup, macOS only)
go install golang.org/x/mobile/cmd/gomobile@latest
go install golang.org/x/mobile/cmd/gobind@latest
gomobile init

# Build for iOS
gomobile build -target=ios -o go_lbm.app .
```

### CI/CD and Releases

This project uses GitHub Actions for automated building and releases:

- **Continuous Integration**: Runs on every push and pull request
  - Code formatting checks
  - Build verification
  - Tests (if available)

- **Multi-platform Builds**: Automated builds for all supported platforms
  - Desktop: Windows, macOS, Linux (x86_64 and ARM64)
  - Mobile: Android and iOS
  - **Note**: Windows ARM64 builds are experimental due to cross-compilation complexity with CGO

- **Automated Releases**: Create a tag starting with `v` to trigger a release
  ```bash
  git tag v1.0.0
  git push origin v1.0.0
  ```

Download pre-built binaries from the [Releases](https://github.com/prasadchandan/go_lbm/releases) page.

### Demo
![Go LBM Demo](https://github.com/prasadchandan/go_lbm/blob/master/repo-assets/demo.gif)

### Credits

 - The code for the solver was ported based on the JavaScript implementation by Dan Schroeder 
   + https://physics.weber.edu/schroeder/fluids/
 - Carlo Barth for the colorMapCreator python script that was used to generate the color maps used in the application
 - [FiraSans Font](https://github.com/mozilla/Fira)
