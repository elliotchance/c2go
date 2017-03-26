/* ldexp example */
#include <stdio.h>      /* printf */
#include <math.h>       /* ldexp */

int main ()
{
  double param, result;
  int n;

  param = 0.95;
  n = 4;
  result = ldexp (param , n);
  printf ("%f * 2^%d = %f\n", param, n, result);
  return 0;
}