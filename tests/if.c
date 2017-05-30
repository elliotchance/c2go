#include <stdio.h>
#include "tests.h"

int main()
{
    plan(2);

    int x = 1;

    // Without else
    if (x == 1)
        pass("%s", "x is equal to one");

    // With else
    if (x != 1)
        fail("%s", "x is not equal to one")
    else
        pass("%s", "x is equal to one");

    done_testing();
}
