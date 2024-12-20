#!/bin/bash

# just running the application with docker-compose because im lazy to do this manually


# Check if the PostgreSQL service is running using Homebrew
if ! brew services list | grep -q "^postgresql.*started"; then
    echo "PostgreSQL service is not running. Starting it now..."
    brew services start postgresql
else
    echo "PostgreSQL service is already running."
fi

# Test via go tests then Run Docker Compose services in detached mode.
echo "Running tests..."
go test -v ./...
if [ $? -ne 0 ]; then
    echo "Tests failed. Exiting..."
    exit 1
fi
echo "Tests passed."
echo "Starting Docker Compose services..."
docker-compose up -d

# Tail logs from all services started by Docker Compose.
echo "Tailing logs from Docker Compose services..."
docker-compose logs -f
