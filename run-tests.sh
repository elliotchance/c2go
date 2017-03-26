#!/bin/bash

function run_test {
    export TEST=$1

    echo $TEST

    # First check that ast2json.py can understand every line of the clang AST.
    clang -Xclang -ast-dump -fsyntax-only $TEST | python ast2json.py > /tmp/0.txt
    if [ $? != 0 ]; then
        cat /tmp/0.txt
        exit 1
    fi

    # Compile with clang
    clang -lm $TEST
    if [ $? != 0 ]; then
        exit 1
    fi
    
    (echo "7" | ./a.out some args) > /tmp/1.txt

    GO_FILES=$(python c2go.py $TEST)
    (echo "7" | go run $GO_FILES some args) > /tmp/2.txt

    diff /tmp/1.txt /tmp/2.txt
    if [ $? != 0 ]; then
        echo "=== out.go"
        cat --number out.go
        exit 1
    fi
}

for TEST in ${1-tests/*.c}; do
    run_test $TEST
done
