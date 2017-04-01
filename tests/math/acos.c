/* acos example */
#include <stdio.h>      /* printf */
#include <math.h>       /* acos */

#define PI 3.14159265

int main ()
{
  double param, result;
  param = 0.5;
  result = acos (param) * 180.0 / PI;
  printf ("The arc cosine of %f is %.3f degrees.\n", param, result);
  return 0;
}