#include <stdio.h>
#include <ctype.h>

int main()
{
   int var1 = '3';
   int var2 = 'm';

   if( isgraph(var1) )
   {
      printf("var1 = |%c| can be printed\n", var1 );
   }
   else
   {
      printf("var1 = |%c| can't be printed\n", var1 );
   }

   if( isgraph(var2) )
   {
      printf("var2 = |%c| can be printed\n", var2 );
   }
   else
   {
      printf("var2 = |%c| can't be printed\n", var2 );
   }

   return(0);
}
