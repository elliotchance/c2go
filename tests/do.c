#include <stdio.h>

int main() {
	int i = 0;

	do {
		printf("do loop %d\n", i);
		i = i + 1;
	} while( i < 4 );
}
