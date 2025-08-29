#!/bin/bash
# Create a directory for storage when using the file system storage in development.
mkdir -p /tmp/howlite-resources-storage/

go install github.com/air-verse/air@latest
cd /workspaces/howlite-resources-v3/src
go mod tidy