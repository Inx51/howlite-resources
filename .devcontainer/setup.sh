#!/bin/bash
mkdir -p /tmp/howlite-resources-storage/
cd /workspaces/howlite-resources/src/
go install github.com/air-verse/air@latest
go mod tidy