// Array examples

#include <stdio.h>
#include "tests.h"

int main()
{
    plan(5);

    int a[3];
    a[0] = 5;
    a[1] = 9;
    a[2] = -13;

    is_eq(a[0], 5);
    is_eq(a[1], 9);
    is_eq(a[2], -13);

    double b[2];
    b[0] = 1.2;
    b[1] = 7; // different type

    is_eq(b[0], 1.2);
    is_eq(b[1], 7.0);
    
    done_testing();
}
