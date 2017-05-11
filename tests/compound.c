#include <stdio.h>

int main() {
	int i = 10;
	float f = 3.14159f;
	double d = 0.0;
	char c = 'A';
	i %= 10; printf("%d\n", i);
	i += 10; printf("%d\n", i);
	i -= 2;  printf("%d\n", i);
	i *= 2;  printf("%d\n", i);
	i /= 4;  printf("%d\n", i);
	i <<= 2; printf("%d\n", i);
	i >>= 2; printf("%d\n", i);
	i ^= 0xCFCF; printf("%d\n", i);
	i |= 0xFFFF; printf("%d\n", i);
	i &= 0x0000; printf("%d\n", i);
	f += 1.0f; d += 1.25f;
	i -= 255l; i += 'A'; c += 11;
	printf("%d %d %f %f\n", i, c, f, d);
}
