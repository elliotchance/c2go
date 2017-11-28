#include <stdio.h>
#include "header.h"

#define ERROR_FUNC error
#include "err.h"
#undef ERROR_FUNC
#define ERROR_FUNC errorANOTHER
#include "err.h"

int main() {
    say_two();
    say_three();
    say_four();

	ERROR_FUNC();
	error();
	errorANOTHER();
    return 0;
}

// The body for the definition (declared in the header). Notice this is declared
// after using the forward declaration above.
void say_two() {
    printf("2");
}
