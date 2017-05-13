// getchar() needs to be in it's own file because it takes input from stdin.

#include <stdio.h>

#define START_TEST(t)         \
    printf("\n--- %s\n", #t); \
    test_##t();

void test_getchar()
{
    int c;
    puts("Enter text. Include a dot ('.') in a sentence to exit:");
    c = getchar();
    putchar(c);
}

int main()
{
    START_TEST(getchar)

    return 0;
}
