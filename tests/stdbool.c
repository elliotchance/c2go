#include <stdio.h>
#include <stdbool.h>
#include "tests.h"

int main()
{
    plan(6);

    bool trueBool = true;
    bool falseBool = false;

    is_true(trueBool);
    is_false(falseBool);

    if (trueBool)
    {
        pass("%s", "true")
    }
    else
    {
        fail("%s", "should not reach here")
    }

    if (!trueBool)
    {
        fail("%s", "should not reach here")
    }
    else
    {
        pass("%s", "true")
    }

    if (falseBool)
    {
        fail("%s", "should not reach here")
    }
    else
    {
        pass("%s", "false")
    }

    if (!falseBool)
    {
        pass("%s", "false")
    }
    else
    {
        fail("%s", "should not reach here")
    }

    done_testing();
}
