#!/bin/sh -eu

sleep 5

echo "Start testing......"
cd /app
go mod download
go install github.com/onsi/ginkgo/v2/ginkgo@v2.1.3
export APP_PATH=/app
ginkgo -r