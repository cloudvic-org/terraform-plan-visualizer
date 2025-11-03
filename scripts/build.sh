#!/bin/bash

# Terraform Plan Visualizer Build Script

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default values
VERSION="dev"
PLATFORMS=""
OUTPUT_DIR="dist"
CLEAN=false
HELP=false

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to show help
show_help() {
    cat << EOF
Terraform Plan Visualizer Build Script

Usage: $0 [OPTIONS]

Options:
    -v, --version VERSION     Set version (default: dev)
    -p, --platforms PLATFORMS
                             Comma-separated list of platforms to build
                             Supported: linux-amd64,linux-arm64,darwin-amd64,darwin-arm64,windows-amd64,windows-arm64
                             Default: all platforms
    -o, --output-dir DIR      Output directory (default: dist)
    -c, --clean              Clean output directory before building
    -h, --help               Show this help message

Examples:
    $0                                    # Build all platforms
    $0 -p linux-amd64,darwin-amd64        # Build specific platforms
    $0 -v 1.0.0 -c                        # Build with version and clean
    $0 -o releases -p windows-amd64      # Build Windows binary to releases/

EOF
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -v|--version)
            VERSION="$2"
            shift 2
            ;;
        -p|--platforms)
            PLATFORMS="$2"
            shift 2
            ;;
        -o|--output-dir)
            OUTPUT_DIR="$2"
            shift 2
            ;;
        -c|--clean)
            CLEAN=true
            shift
            ;;
        -h|--help)
            HELP=true
            shift
            ;;
        *)
            print_error "Unknown option: $1"
            show_help
            exit 1
            ;;
    esac
done

# Show help if requested
if [ "$HELP" = true ]; then
    show_help
    exit 0
fi

# Default platforms if none specified
if [ -z "$PLATFORMS" ]; then
    PLATFORMS="linux-amd64,linux-arm64,darwin-amd64,darwin-arm64,windows-amd64,windows-arm64"
fi

# Clean output directory if requested
if [ "$CLEAN" = true ]; then
    print_status "Cleaning output directory: $OUTPUT_DIR"
    rm -rf "$OUTPUT_DIR"
fi

# Create output directory
print_status "Creating output directory: $OUTPUT_DIR"
mkdir -p "$OUTPUT_DIR"

# Get build information
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

print_status "Build Information:"
print_status "  Version: $VERSION"
print_status "  Build Time: $BUILD_TIME"
print_status "  Git Commit: $GIT_COMMIT"
print_status "  Platforms: $PLATFORMS"

# Build flags
LDFLAGS="-ldflags \"-X main.Version=$VERSION -X main.BuildTime=$BUILD_TIME -X main.GitCommit=$GIT_COMMIT\""

# Function to build for a specific platform
build_platform() {
    local platform=$1
    local os=$(echo $platform | cut -d'-' -f1)
    local arch=$(echo $platform | cut -d'-' -f2)
    
    print_status "Building for $platform ($os/$arch)..."
    
    local binary_name="terraform-plan-visualizer"
    if [ "$os" = "windows" ]; then
        binary_name="${binary_name}.exe"
    fi
    
    local output_file="$OUTPUT_DIR/terraform-plan-visualizer-$platform"
    if [ "$os" = "windows" ]; then
        output_file="${output_file}.exe"
    fi
    
    # Build the binary
    GOOS=$os GOARCH=$arch go build $LDFLAGS -o "$output_file" .
    
    if [ $? -eq 0 ]; then
        print_success "Built $platform successfully: $output_file"
        
        # Show file size
        local size=$(du -h "$output_file" | cut -f1)
        print_status "  Size: $size"
    else
        print_error "Failed to build $platform"
        return 1
    fi
}

# Check if Go is installed
if ! command -v go &> /dev/null; then
    print_error "Go is not installed or not in PATH"
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | cut -d' ' -f3 | sed 's/go//')
print_status "Using Go version: $GO_VERSION"

# Build for each platform
print_status "Starting build process..."
BUILD_COUNT=0
FAILED_BUILDS=0

IFS=',' read -ra PLATFORM_ARRAY <<< "$PLATFORMS"
for platform in "${PLATFORM_ARRAY[@]}"; do
    platform=$(echo $platform | xargs) # Trim whitespace
    
    if build_platform "$platform"; then
        ((BUILD_COUNT++))
    else
        ((FAILED_BUILDS++))
    fi
done

# Summary
print_status "Build Summary:"
print_status "  Successful builds: $BUILD_COUNT"
if [ $FAILED_BUILDS -gt 0 ]; then
    print_warning "  Failed builds: $FAILED_BUILDS"
fi

# List all built binaries
print_status "Built binaries:"
ls -la "$OUTPUT_DIR"/* 2>/dev/null || print_warning "No binaries found in $OUTPUT_DIR"

if [ $FAILED_BUILDS -eq 0 ]; then
    print_success "All builds completed successfully!"
    exit 0
else
    print_error "Some builds failed!"
    exit 1
fi
