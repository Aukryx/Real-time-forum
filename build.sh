#!/bin/bash
set -e

# Go to the source directory
cd cmd/golang-server-layout

# Build the app
go build -o ../../app

# Return to root
cd ../..

# Make executable
chmod +x app

# Create directory structure
echo "Creating directory structure..."
mkdir -p web/templates

# Copy template files to ensure they exist in both potential locations
echo "Copying template files..."
cp -r cmd/golang-server-layout/../../web/templates/* web/templates/ || echo "Warning: Template copy failed, may not exist at source"

# Copy static files
echo "Copying static files..."
mkdir -p web/static
cp -r cmd/golang-server-layout/../../web/static/* web/static/ || echo "Warning: Static files copy failed, may not exist at source"

# Debug information
echo "Final web directory structure:"
find web -type f | head -n 10
echo "Total files: $(find web -type f | wc -l)"

echo "Build complete!"