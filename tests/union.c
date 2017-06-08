// Tests for unions.

#include <stdio.h>

union programming
{
    float constant;
    char *pointer;
};

void f(union programming *addr)
{
    char *s = "Show string member.";

    addr->constant += 4.56;
    addr->constant++;
    printf("%f\n", addr->constant);

    addr->pointer = s;
    printf("%s\n", addr->pointer);
}

int main()
{
    union programming variable;
    char *s = "Programming in Software Development.";

    /*union {
        float constant;
        char *pointer;
    } local_union = */

    variable.constant = 1.23;
    printf("%f\n", variable.constant);

    variable.pointer = s;
    printf("%s\n", variable.pointer);

    f(&variable);

    return 0;
}
