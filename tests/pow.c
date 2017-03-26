/* pow example */
#include <stdio.h>      /* printf */
#include <math.h>       /* pow */

int main ()
{
  printf ("7 ^ 3 = %f\n", pow (7.0, 3.0) );
  printf ("4.73 ^ 12 = %f\n", pow (4.73, 12.0) );
  printf ("32.01 ^ 1.54 = %f\n", pow (32.01, 1.54) );
  return 0;
}