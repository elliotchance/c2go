#include <stdio.h>

int main()
{
   int i;
   for (i = 0; i < 10; i++)
	printf("%d\n", i);

   for (;;)
        printf("infinite loop");
}
