// This file contains tests for the system arguments (argv).

#include <stdio.h>
#include "tests.h"

int main(int argc, const char **argv)
{
    plan(3);

    is_eq(argc, 3);

    // We cannot compare the zeroth argument becuase it will be different for C
    // and Go.
    // is_streq(argv[0], "build/go.out");

    is_streq(argv[1], "some");
    is_streq(argv[2], "args");

    done_testing();
}
