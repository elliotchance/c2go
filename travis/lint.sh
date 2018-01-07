#!/bin/bash

set -e

CLANG_FORMAT="clang-format-3.5"

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

# Install clang-format
if [ "$TRAVIS_OS_NAME" = "linux" ]; then
	sudo apt-get update
	sudo apt-cache search clang
	sudo apt-get install -f -y --force-yes $CLANG_FORMAT
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
		echo "Diff:"
		diff $C_FILE /tmp/out
    	exit 1
	fi
done
