#include "tests.h"
#include <stdbool.h>
#include <stdio.h>

_Bool f(_Bool b)
{
    return b;
}

int main()
{
    plan(12);

    bool trueBool = true;
    bool falseBool = false;

    is_true(trueBool);
    is_false(falseBool);

    if (trueBool) {
        pass("%s", "true")
    } else {
        fail("%s", "should not reach here")
    }

    if (!trueBool) {
        fail("%s", "should not reach here")
    } else {
        pass("%s", "true")
    }

    if (falseBool) {
        fail("%s", "should not reach here")
    } else {
        pass("%s", "false")
    }

    if (!falseBool) {
        pass("%s", "false")
    } else {
        fail("%s", "should not reach here")
    }

    _Bool var = true;
    if (var) {
        pass("%s", "ok")
    } else {
        fail("%s", "should not reach here")
    }

    var = true;
    if (var - var) {
        fail("%s", "should not reach here")
    } else {
        pass("%s", "ok")
    }

    var = true;
    if (var - var == false) {
        pass("%s", "ok")
    } else {
        fail("%s", "should not reach here")
    }

    _Bool b = 0; // false
    if (b) {
        b++;
    }
    if (b == false) // b = 0
    {
        pass("%s", "ok")
    }

    _Bool c = f(b);
    b = b + c;
    if (b == false) {
        pass("%s", "ok")
    }
    int i = (int)(b);
    if (i == 0) {
        pass("%s", "ok")
    }

    done_testing();
}
