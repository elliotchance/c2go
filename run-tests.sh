#!/bin/bash

for TEST in $(ls -1 tests); do
    echo $TEST

    # First check that ast2json.py can understand every line of the clang AST.
    clang -Xclang -ast-dump -fsyntax-only tests/argv.c | python ast2json.py > /tmp/0.txt
    if [ $? != 0 ]; then
        cat /tmp/0.txt
        exit 1
    fi

    clang tests/$TEST
    (echo "7" | ./a.out some args) > /tmp/1.txt

    python c2go.py tests/$TEST > out.go
    (echo "7" | go run functions.go out.go some args) > /tmp/2.txt

    diff /tmp/1.txt /tmp/2.txt
    if [ $? != 0 ]; then
        exit 1
    fi
done
