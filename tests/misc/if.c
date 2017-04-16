#include <stdio.h>
 
int main()
{
   int x = 1;

   // Without else
   if ( x == 1 )
      printf("x is equal to one.\n");

   // With else
   if ( x != 1 )
      printf("x is not equal to one.\n");
   else
      printf("x is equal to one.\n");
 
   return 0;
}
