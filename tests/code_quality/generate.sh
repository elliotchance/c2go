#!/bin/bash

# These steps are from the README to verify it can be installed and run as
# documented.
go build

# Generate code quality Go code
FILES='tests/code_quality/*.c'
for file in $FILES
do
  echo "Processing $file file..."
  c2go transpile -o="$file.go" $file
done
