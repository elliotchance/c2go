#include <stdio.h>

// This header file is included by multiple C files. Make sure the preprocessor
// is working correctly by not declaring duplicates.
#ifndef HEADER
#define HEADER

// Headers usually do not contain whole functions, but we want to make sure this
// is still OK to do.
void say_four() {
    printf("4");
}

// Forward-declared prototypes that are defined in one of our other C files.
void say_two();   // main1.c
void say_three(); // main2.c

#endif /* HEADER */
