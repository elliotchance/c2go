#include <assert.h>
#include <stdio.h>
#include <stdlib.h>

void test_malloc1()
{
    int i, n;
    char *buffer;

    printf("How long do you want the string? ");
    scanf("%d", &i);

    buffer = (char *)malloc(i + 1);
    if (buffer == NULL)
        return;

    for (n = 0; n < i; n++)
        buffer[n] = i % 26 + 'a';
    buffer[i] = '\0';

    printf("Random string: %s\n", buffer);
    free(buffer);
}

void test_malloc2()
{
    int *p;
    p = (int *)malloc(sizeof(int));
    *p = 5;
}

// Mix around all the types to make sure it is still actually allocating the
// correct size.
void test_malloc3()
{
    assert(sizeof(int) == 4);
    assert(sizeof(double) == 8);

    // 10 ints, should be 5 doubles. Also use a bad cast to make sure that it
    // doesn't interfere with the types.
    double *d;
    d = (char *)malloc(sizeof(int) * 10);

    // We can't test how much memory has been allocated by Go, but we can test
    // that there are 5 items.
    *d = 123;
    d[4] = 456;

    printf("%f %f", d[0], d[4]);
}

int main()
{
    test_malloc1();
    test_malloc2();
    test_malloc3();

    return 0;
}
