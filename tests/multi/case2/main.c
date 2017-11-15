#include"head.h"

#ifndef HEADER
#define HEADER

#include <stdio.h>

void say_four(){
	printf("4");
}
void say_two(){
	printf("2");
}
#endif /* HEADER */

int main(){
	say_four();
	say_two();
	return 0;
}
