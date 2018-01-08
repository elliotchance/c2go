// This file contains tests for the system arguments (argv).

#include "tests.h"
#include <stdio.h>

int main(int argc, const char** argv)
{
    plan(2);

    // When this file is converted to go it is run through "go test" that needs
    // some extra arguments before the standard C arguments. We need to adjust
    // an offset so that the C program and the Go program read the same index
    // for the first index of the real arguments.
    int offset = 0;

    // More than three arguments means it must be run under "go test". If not
    // the assertion immediately below will fail.
    if (argc > 3) {
        offset = 3;
    }

    // We cannot compare the zeroth argument because it will be different for C
    // and Go.
    // is_streq(argv[0], "build/go.out");

    is_streq(argv[1 + offset], "some");
    is_streq(argv[2 + offset], "args");

    done_testing();
}
