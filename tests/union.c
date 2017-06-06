// Tests for unions.

#include <stdio.h>
#include "tests.h"

union programming
{
    float constant;
    char *pointer;
};

void f(union programming *addr)
{
    char *s = "Show string member.";
    float v = 1.23+4.56+1.;

    addr->constant += 4.56;
    addr->constant++;
    is_eq(addr->constant, v);

    addr->pointer = s;
    is_streq(addr->pointer, "Show string member.");
}

int main()
{
    plan(4);

    union programming variable;
    char *s = "Programming in Software Development.";

    variable.pointer = s;
    is_streq(variable.pointer, "Programming in Software Development.");

    variable.constant = 1.23;
    is_eq(variable.constant, 1.23);

    f(&variable);

    done_testing();

    return 0;
}
