#!/bin/bash
# Build script for atuin-ui
# This script builds the frontend, embeds it in the backend, and compiles for multiple platforms

# Stop on error
set -e

# Configuration
FRONTEND_DIR="frontend"
BACKEND_DIR="backend"
DIST_DIR="dist"
FRONTEND_DIST_DIR="frontend/dist"
BACKEND_EMBED_DIR="backend/frontend_dist"
APP_NAME="atuin-web-ui"

# Ensure dist directory exists
mkdir -p "$DIST_DIR"
echo "Created dist directory"

# Build frontend
echo "Building frontend..."
cd "$FRONTEND_DIR"

# Install dependencies if needed
if [ ! -d "node_modules" ]; then
    echo "Installing frontend dependencies..."
    npm install
fi

# Build frontend
npm run build
echo "Frontend built successfully"
cd ..

# Prepare backend embed directory
echo "Preparing backend embed directory..."
if [ -d "$BACKEND_EMBED_DIR" ]; then
    rm -rf "$BACKEND_EMBED_DIR"
fi
mkdir -p "$BACKEND_EMBED_DIR"

# Copy frontend build to backend embed directory
echo "Copying frontend build to backend embed directory..."
cp -r "$FRONTEND_DIST_DIR"/* "$BACKEND_EMBED_DIR"

# Ensure the directory structure exists
echo "Ensuring directory structure exists..."
directories=("css" "js" "img" "fonts")
for dir in "${directories[@]}"; do
    path="$BACKEND_EMBED_DIR/$dir"
    if [ ! -d "$path" ]; then
        mkdir -p "$path"
        echo "Created directory: $path"
    fi
done

# Build backend for different platforms
echo "Building backend for multiple platforms..."

# Get current OS
CURRENT_OS=$(uname -s | tr '[:upper:]' '[:lower:]')
if [[ "$CURRENT_OS" == *"darwin"* ]]; then
    CURRENT_OS="darwin"
elif [[ "$CURRENT_OS" == *"linux"* ]]; then
    CURRENT_OS="linux"
elif [[ "$CURRENT_OS" == *"mingw"* ]] || [[ "$CURRENT_OS" == *"msys"* ]] || [[ "$CURRENT_OS" == *"cygwin"* ]]; then
    CURRENT_OS="windows"
fi

# Get current architecture
CURRENT_ARCH=$(uname -m)
if [[ "$CURRENT_ARCH" == "x86_64" ]]; then
    CURRENT_ARCH="amd64"
elif [[ "$CURRENT_ARCH" == "arm64" ]] || [[ "$CURRENT_ARCH" == "aarch64" ]]; then
    CURRENT_ARCH="arm64"
fi

echo "Current platform: $CURRENT_OS/$CURRENT_ARCH"

# Define build targets - now we can build for all platforms
TARGETS=(
    "windows amd64 .exe"
    "darwin amd64 ''"
    "darwin arm64 ''"
    "linux amd64 ''"
)

echo "Building for all platforms"

# Build for each target
for target in "${TARGETS[@]}"; do
    read -r goos goarch suffix <<< "$target"
    output_name="${APP_NAME}-${goos}-${goarch}${suffix}"
    output_path="${DIST_DIR}/${output_name}"

    echo "Building for $goos/$goarch..."

    cd "$BACKEND_DIR"

    GOOS="$goos" GOARCH="$goarch" go build -o "../$output_path" -ldflags="-s -w" .
    echo "Built $output_path successfully"
    cd ..
done

echo "Build completed successfully!"
echo "Binaries are available in the $DIST_DIR directory"
