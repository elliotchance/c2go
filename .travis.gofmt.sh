#!/bin/bash

if [ -n "$(gofmt -l .)" ]; then
    echo "Go code is not properly formatted. Use 'gofmt'."
    gofmt -d .
    exit 1
fi
