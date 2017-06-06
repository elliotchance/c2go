#include <stdio.h>
#include "tests.h"

int main()
{
    plan(4);

    int value = 1;

    while (value <= 3)
    {
        pass("value is %d", value);
        value++;
    }

    // continue
    value = 0;
    while (value < 3)
    {
        value++;
        if (value < 3)
            continue;
        pass("%d", value);
    }

    done_testing();
}
