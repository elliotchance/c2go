#!/bin/bash

set -e

# clang-format version
CLANG_FORMAT="clang-format"

# Arguments menu
echo "    -r rewrite C test files in according to code-style"
if [ "$1" == "-r" ]; then
	C_TEST_FILES=`ls ./tests/*.c`
	for C_FILE in $C_TEST_FILES
	do
		echo "Formatting file '$C_FILE' ..."
		eval "$CLANG_FORMAT -style=WebKit -i $C_FILE"
	done
fi

# Check go fmt first
if [ -n "$(gofmt -l .)" ]; then
    echo "Go code is not properly formatted. Use 'gofmt'."
    gofmt -d .
    exit 1
fi

# Version of clang-format
echo "Version of clang-format:"
eval "$CLANG_FORMAT -version"

# Check clang-format for C test source files
C_TEST_FILES=`ls ./tests/*.c`
for C_FILE in $C_TEST_FILES
do
	eval "$CLANG_FORMAT -style=WebKit $C_FILE > /tmp/out"
	if [ -n "$(diff $C_FILE /tmp/out)" ]; then
    	echo "C test code '$C_FILE' is not properly formatted. Use '$CLANG_FORMAT -style=WebKit'."
	fi
done
