// Tests for unions.

#include <stdio.h>

union programming
{
    float constant;
    char *pointer;
};

int main()
{
    union programming variable;
    char *s = "Programming in Software Development.";

    variable.constant = 1.23;
    printf("%f\n", variable.constant);

    variable.pointer = s;
    printf("%s\n", variable.pointer);

    return 0;
}
