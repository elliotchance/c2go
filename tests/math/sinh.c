/* sinh example */
#include <stdio.h>      /* printf */
#include <math.h>       /* sinh, log */

int main ()
{
  double param, result;
  param = log(2.0);
  result = sinh (param);
  printf ("The hyperbolic sine of %f is %f.\n", param, result );
  return 0;
}