#include <stdio.h>

#ifndef HEADER
#define HEADER

void say_four(){
	printf("4");
}
void say_two(){
	printf("2");
}
#endif /* HEADER */

#ifndef HEADER2
#define HEADER2

#include <stdio.h>

#endif
