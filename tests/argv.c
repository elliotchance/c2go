// This file contains tests for the system arguments (argv).
//
// TODO: 'argc' and 'argv' are hard-coded.

#include <stdio.h>

int main(int argc, char *argv[])
{
    int c;

    printf("Number of command line arguments passed: %d\n", argc);

    for (c = 1; c < argc; c++)
        printf("%d. Command line argument passed is %s\n", c + 1, argv[c]);

    return 0;
}
