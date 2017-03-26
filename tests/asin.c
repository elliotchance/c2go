/* asin example */
#include <stdio.h>      /* printf */
#include <math.h>       /* asin */

#define PI 3.14159265

int main ()
{
  double param, result;
  param = 0.5;
  result = asin (param) * 180.0 / PI;
  printf ("The arc sine of %f is %.3f degrees\n", param, result);
  return 0;
}