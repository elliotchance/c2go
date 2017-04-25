// This program actually still works without including stdio.h but it should be
// here for consistency.

#include <stdio.h>

void test_putchar()
{
    char c;
    for (c = 'A'; c <= 'Z'; c++)
        putchar(c);
}

void test_puts()
{
    puts("c2go");
}

void test_printf()
{
    printf("Hello World\n");
}

void test_remove()
{
    if (remove("myfile.txt") != 0)
        puts("Error deleting file");
    else
        puts("File successfully deleted");
}

int main()
{
    test_putchar();
    test_puts();
    test_printf();
    test_remove();

    return 0;
}
