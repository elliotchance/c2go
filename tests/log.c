/* log example */
#include <stdio.h>      /* printf */
#include <math.h>       /* log */

int main ()
{
  double param, result;
  param = 5.5;
  result = log (param);
  printf ("log(%f) = %f\n", param, result );
  return 0;
}