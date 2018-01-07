#include "tests.h"
#include <stdio.h>

#define START_TEST(t) \
    diag(#t);         \
    test_##t();

void test_notint()
{
    int i = 0;
    if (!i) {
        pass("good");
    } else {
        fail("fail");
    }

    i = 123;
    if (!i) {
        fail("fail");
    } else {
        pass("good");
    }
}

void test_notptr()
{
    FILE* fp = NULL;
    if (!fp) {
        pass("good");
    } else {
        fail("fail");
    }

    fp = stdin;
    if (!fp) {
        fail("fail");
    } else {
        pass("good");
    }
}

int main()
{
    plan(4);

    START_TEST(notint)
    START_TEST(notptr)

    done_testing();
}
