// Tests for unions.

#include <stdio.h>
#include "tests.h"

union programming
{
    float constant;
    char *pointer;
};

union programming init_var()
{
    union programming variable;
    char *s = "Programming in Software Development.";

    variable.pointer = s;
    is_streq(variable.pointer, "Programming in Software Development.");

    variable.constant = 1.23;
    is_eq(variable.constant, 1.23);

    return variable;
}

void pass_by_ref(union programming *addr)
{
    char *s = "Show string member.";
    float v = 1.23+4.56;

    addr->constant += 4.56;
    is_eq(addr->constant, v);

    addr->pointer = s;
    is_streq(addr->pointer, "Show string member.");
}

void var_by_val(union programming value)
{
    value.constant++;

    is_eq(value.constant, 2.23);
}

int main()
{
    plan(5);

    union programming variable;

    variable = init_var();
    var_by_val(variable);
    pass_by_ref(&variable);

    done_testing();
}
