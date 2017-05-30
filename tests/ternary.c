#include <stdio.h>

int main() {

	int a = 'a' == 65 ? 10 : 100;
	float b = 10 == 10 ? 1.0 : 2.0;
	char* c = 'x' == 5 ? "one" : "two";
	char d = a == 100 ? 'x' : 1;
	printf("%d %f %s %d\n", a, b, c, d);
	printf("%d %d %d\n", 0?1:0, NULL?1:0, 'x'?1:0);

	a = a == 10 ? b == 1.0 ? 1 : 2 : 2;

	if( a == (a == 2 ? 5 : 10) ) {
		printf("%d\n", a);
	}

}
