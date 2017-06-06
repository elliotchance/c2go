// Tests for structures.

#include <stdio.h>

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
    struct programming variable;
    char *s = "Programming in Software Development.";

    variable.constant = 1.23;
    variable.pointer = s;

    printf("%f\n", variable.constant);
    printf("%s\n", variable.pointer);

    f(&variable);

    return 0;
}
