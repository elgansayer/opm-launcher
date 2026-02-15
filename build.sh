#!/bin/bash

# Configuration
BINARY_NAME="opm-launcher"
BUILD_DIR="bin"

# Ensure build directory exists
mkdir -p $BUILD_DIR

# Define platforms to build for: GOOS/GOARCH/EXTENSION
platforms=(
    "linux/amd64"
    "linux/arm64"
    "windows/amd64/.exe"
    "darwin/amd64"
    "darwin/arm64"
)

echo "Starting cross-platform build..."

for platform in "${platforms[@]}"; do
    # Split the platform string
    IFS="/" read -r -a parts <<< "$platform"
    GOOS=${parts[0]}
    GOARCH=${parts[1]}
    EXT=${parts[2]}
    
    OUTPUT_NAME="${BINARY_NAME}-${GOOS}-${GOARCH}${EXT}"
    
    echo "  -> Building $OUTPUT_NAME..."
    
    # Execute the build
    env GOOS=$GOOS GOARCH=$GOARCH go build -o $BUILD_DIR/$OUTPUT_NAME main.go
    
    if [ $? -ne 0 ]; then
        echo "ERROR: Failed to build for $GOOS/$GOARCH"
        exit 1
    fi
done

echo "Build complete! Artifacts are in the '$BUILD_DIR' directory."
ls -lh $BUILD_DIR
