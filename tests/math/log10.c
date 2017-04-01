/* log10 example */
#include <stdio.h>      /* printf */
#include <math.h>       /* log10 */

int main ()
{
  double param, result;
  param = 1000.0;
  result = log10 (param);
  printf ("log10(%f) = %f\n", param, result );
  return 0;
}