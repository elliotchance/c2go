#!/bin/bash

set -e

# Check go fmt first
if [ -n "$(gofmt -l .)" ]; then
    echo "Go code is not properly formatted. Use 'gofmt'."
    gofmt -d .
    exit 1
fi
