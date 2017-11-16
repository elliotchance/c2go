#include <stdio.h>
#include "header.h"

int main() {
    say_two();
    say_three();
    say_four();

    return 0;
}

// The body for the definition (declared in the header). Notice this is declared
// after using the forward declaration above.
void say_two() {
    printf("2");
}
