// This file contains tests for assert.h.
//
// Note: assert() failures print directly. The output of this program will not
// be valid TAP.

#include "tests.h"
#include <assert.h>
#include <stdio.h>

void print_number(int* myInt)
{
    assert(myInt != NULL);
    printf("%d\n", *myInt);
}

int main()
{
    plan(0);

    int a = 10;
    int* b = NULL;
    int* c = NULL;

    b = &a;

    print_number(b);
    print_number(c);

    fail("%s", "It shouldn't make it to here!");
    done_testing();
}
