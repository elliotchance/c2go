#include <stdio.h>
int main(){
	// Trigraph tests require the `-trigraph` clang option, so it cannot be amount other general tests.
	printf("??=");
	return 0;
}
