#!/bin/bash

go build

# Generate code quality Go code
FILES='tests/code_quality/*.c'
for file in $FILES
do
  echo "Processing $file file..."
  c2go transpile -o="$file.txt" -p="code_quality" $file
done
