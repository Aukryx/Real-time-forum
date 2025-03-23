#!/bin/bash
set -e

# Set PORT from environment variable
export PORT="${PORT:-8080}"

# Enhanced diagnostic information
echo "==== DEPLOYMENT DIAGNOSTICS ===="
echo "Starting server on port: $PORT"
echo "Current directory: $(pwd)"
echo "Files in current directory: $(ls -la)"

# Check for template files specifically
echo "==== TEMPLATE DIAGNOSTICS ===="
echo "Checking for templates in ../../web/templates:"
ls -la ../../web/templates 2>/dev/null || echo "Directory not found"

echo "Checking for templates in ./web/templates:"
ls -la ./web/templates 2>/dev/null || echo "Directory not found"

echo "Checking for templates in web/templates:"
ls -la web/templates 2>/dev/null || echo "Directory not found"

# Check for index.html specifically
echo "Searching for index.html:"
find . -name "index.html" 2>/dev/null || echo "No index.html found"

# Final directory structure
echo "==== WEB DIRECTORY STRUCTURE ===="
find web -type f 2>/dev/null | sort || echo "No web directory found"

echo "==== STARTING APPLICATION ===="
# Run the app
./app