#!/bin/bash

echo "===== SENTINEL Release Builder ====="

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go not found. Make sure Go is installed and in your PATH."
    exit 1
fi

# Check if GoReleaser is installed
if ! command -v goreleaser &> /dev/null; then
    echo "GoReleaser not found. Installing GoReleaser..."
    go install github.com/goreleaser/goreleaser@latest
    if [ $? -ne 0 ]; then
        echo "Error: Failed to install GoReleaser."
        exit 1
    fi
    echo "GoReleaser installed successfully."
fi

# Clean dist folder if it exists
if [ -d "dist" ]; then
    echo "Cleaning dist folder..."
    rm -rf dist
fi

# Run GoReleaser in snapshot mode (without uploading to GitHub)
echo "Creating local release..."
goreleaser release --snapshot --clean

if [ $? -ne 0 ]; then
    echo "Error: Failed to create release."
    exit 1
fi

echo
echo "===== Release created successfully! ====="
echo "Binaries available in dist/ folder"
echo

# Display list of generated files
echo "Generated files:"
ls -la dist/