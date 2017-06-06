// Tests for structures.

#include <stdio.h>
#include "tests.h"

struct programming
{
    float constant;
    char *pointer;
};

void f(struct programming *addr)
{
    char *s = "Show string member.";
    float v = 1.23+4.56+1.;

    addr->constant += 4.56;
    addr->constant++;
    addr->pointer = s;

    is_eq(addr->constant, v);
    is_streq(addr->pointer, "Show string member.");
}

int main()
{
    plan(4);

    struct programming variable;
    char *s = "Programming in Software Development.";

    variable.constant = 1.23;
    variable.pointer = s;

    is_eq(variable.constant, 1.23);
    is_streq(variable.pointer, "Programming in Software Development.");

    f(&variable);

    done_testing();

    return 0;
}
