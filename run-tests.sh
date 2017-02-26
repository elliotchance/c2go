#!/bin/bash

for TEST in $(ls -1 tests); do
    echo $TEST

    clang tests/$TEST
    (echo "7" | ./a.out) > /tmp/1.txt

    python c2go.py tests/$TEST > out.go
    (echo "7" | go run functions.go out.go) > /tmp/2.txt

    diff /tmp/1.txt /tmp/2.txt
    if [ $? != 0 ]; then
        exit
    fi
done
