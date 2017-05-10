#include <stdio.h>

int main() {
	int i = 0;

	do {
		printf("do loop %d\n", i);
		i = i + 1;
	} while( i < 4 );

	// continue
	i = 0;
	do {
		i++;
		if(i < 3) continue;
		printf("%d\n", i);
	} while(i < 3);

}
