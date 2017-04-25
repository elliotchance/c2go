// Tests for structures.

#include <stdio.h>

struct programming
{
    float constant;
    char *pointer;
};

int main()
{
    struct programming variable;
    char *s = "Programming in Software Development.";

    variable.constant = 1.23;
    variable.pointer = s;

    printf("%f\n", variable.constant);
    printf("%s\n", variable.pointer);

    return 0;
}
