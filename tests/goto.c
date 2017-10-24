#include <stdio.h>
#include <time.h>

#include "tests.h"

#define START_TEST(t) \
    diag(#t);         \
    test_##t();

void test_goto()
{
    int i = 0;
    
    mylabel: i++;
    
    if (i == 1) {
        goto mylabel;
    }

    is_eq(i, 2);
}

void test_goto_stmt()
{
    int i = 0, j = 0;
    
    mylabel: for (j=0; j<5; j++) i++;
    
    if (i < 15) {
        goto mylabel;
    }

    is_eq(i, 15);
}

int main()
{
    plan(2);

    START_TEST(goto)
    START_TEST(goto_stmt)
    
    done_testing();
}
