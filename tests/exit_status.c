#include "tests.h"
#include <stdio.h>

int main()
{
    plan(0);

    // There is no done_testing() becuase we want to return an error code without
    // checks that fail.
    return 1;
}
