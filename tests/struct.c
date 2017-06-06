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

    addr->constant += 4.56;
    addr->constant++;
    addr->pointer = s;

    printf("%f\n", addr->constant);
    printf("%s\n", addr->pointer);
}

int main()
{
    plan(2);

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
