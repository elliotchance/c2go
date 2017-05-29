// This file contains tests for the system arguments (argv).

#include <stdio.h>
#include "tests.h"

int main(int argc, const char **argv)
{
    plan(4);

    is_eq(argc, 3);

    is_streq(argv[0], "build/go.out");
    is_streq(argv[1], "some");
    is_streq(argv[2], "args");

    done_testing();
}
