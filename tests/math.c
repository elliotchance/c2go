// Test for math.h.

#include <stdio.h>
#include <math.h>
#include "tests.h"

#define PI 3.14159265
#define IS_NAN -2147483648

unsigned long long ullmax = 18446744073709551615ull;

void test_asin()
{
  double param, result;
  param = 0.5;
  result = asin(param) * 180.0 / PI;
  printf("The arc sine of %f is %.3f degrees\n", param, result);
}

void test_atan()
{
  double param, result;
  param = 1.0;
  result = atan(param) * 180 / PI;
  printf("The arc tangent of %f is %.3f degrees\n", param, result);
}

void test_atan2()
{
  double x, y, result;
  x = -10.0;
  y = 10.0;
  result = atan2(y, x) * 180 / PI;
  printf("The arc tangent for (x=%f, y=%f) is %.3f degrees\n", x, y, result);
}

void test_ceil()
{
  printf("ceil of 2.3 is %.1f\n", ceil(2.3));
  printf("ceil of 3.8 is %.1f\n", ceil(3.8));
  printf("ceil of -2.3 is %.1f\n", ceil(-2.3));
  printf("ceil of -3.8 is %.1f\n", ceil(-3.8));
}

void test_cos()
{
  double param, result;
  param = 60.0;
  result = cos(param * PI / 180.0);
  printf("The cosine of %f degrees is %f.\n", param, result);
}

void test_cosh()
{
  double param, result;
  param = log(2.0);
  result = cosh(param);
  printf("The hyperbolic cosine of %f is %f.\n", param, result);
}

void test_exp()
{
  double param, result;
  param = 5.0;
  result = exp(param);
  printf("The exponential value of %f is %f.\n", param, result);
}

void test_fabs()
{
  printf("The absolute value of 3.1416 is %f\n", fabs(3.1416));
  printf("The absolute value of -10.6 is %f\n", fabs(-10.6));
}

void test_floor()
{
  printf("floor of 2.3 is %.1f\n", floor(2.3));
  printf("floor of 3.8 is %.1f\n", floor(3.8));
  printf("floor of -2.3 is %.1f\n", floor(-2.3));
  printf("floor of -3.8 is %.1f\n", floor(-3.8));
}

void test_fmod()
{
  printf("fmod of 5.3 / 2 is %f\n", fmod(5.3, 2));
  printf("fmod of 18.5 / 4.2 is %f\n", fmod(18.5, 4.2));
}

void test_ldexp()
{
  double param, result;
  int n;

  param = 0.95;
  n = 4;
  result = ldexp(param, n);
  printf("%f * 2^%d = %f\n", param, n, result);
}

void test_log()
{
  double param, result;
  param = 5.5;
  result = log(param);
  printf("log(%f) = %f\n", param, result);
}

void test_log10()
{
  double param, result;
  param = 1000.0;
  result = log10(param);
  printf("log10(%f) = %f\n", param, result);
}

void test_pow()
{
  printf("7 ^ 3 = %f\n", pow(7.0, 3.0));
  printf("4.73 ^ 12 = %f\n", pow(4.73, 12.0));
  printf("32.01 ^ 1.54 = %f\n", pow(32.01, 1.54));
}

void test_sin()
{
  double param, result;
  param = 30.0;
  result = sin(param * PI / 180);
  printf("The sine of %f degrees is %f.\n", param, result);
}

void test_sinh()
{
  double param, result;
  param = log(2.0);
  result = sinh(param);
  printf("The hyperbolic sine of %f is %f.\n", param, result);
}

void test_sqrt()
{
  double param, result;
  param = 1024.0;
  result = sqrt(param);
  printf("sqrt(%f) = %f\n", param, result);
}

void test_tan()
{
  double param, result;
  param = 45.0;
  result = tan(param * PI / 180.0);
  printf("The tangent of %f degrees is %f.\n", param, result);
}

void test_tanh()
{
  double param, result;
  param = log(2.0);
  result = tanh(param);
  printf("The hyperbolic tangent of %f is %f.\n", param, result);
}

// Test if x == -0.0.
int isnegzero(double x)
{
  return (x * -0.0) == 0.0 && signbit(x);
}

int main()
{
  plan(122);

  // Test constants
  diag("constants");
  eq_ok(M_E, 2.718282);
  eq_ok(M_LOG2E, 1.442695);
  eq_ok(M_LOG10E, 0.434294);
  eq_ok(M_LN2, 0.693147);
  eq_ok(M_LN10, 2.302585);
  eq_ok(M_PI, 3.141593);
  eq_ok(M_PI_2, 1.570796);
  eq_ok(M_PI_4, 0.785398);
  eq_ok(M_1_PI, 0.318310);
  eq_ok(M_2_PI, 0.636620);
  eq_ok(M_2_SQRTPI, 1.128379);
  eq_ok(M_SQRT2, 1.414214);
  eq_ok(M_SQRT1_2, 0.707107);

  // Each of the tests are against these values:
  //
  // * Simple: 0, 1, -1, 0.5
  // * Large and small: 1.23e300, -1.23e-300
  // * Constants: M_PI, M_E
  // * Special: INFINITY, -INFINITY, NAN

  diag("acos");

  eq_ok(acos(0), 1.570796);
  eq_ok(acos(1), 0);
  eq_ok(acos(-1), M_PI);
  eq_ok(acos(0.5), 1.047198);

  ok(isnan(acos(1.23e300)));
  eq_ok(acos(-1.23e-300), 1.570796);

  ok(isnan(acos(M_PI)));
  ok(isnan(acos(M_E)));

  ok(isnan(acos(INFINITY)));
  ok(isnan(acos(-INFINITY)));
  ok(isnan(acos(NAN)));

  diag("asin");

  eq_ok(asin(0), 0);
  eq_ok(asin(1), 1.570796);
  eq_ok(asin(-1), -1.570796);
  eq_ok(asin(0.5), 0.523599);

  ok(isnan(asin(1.23e300)));
  ok(isnegzero(asin(-1.23e-300)));

  ok(isnan(asin(M_PI)));
  ok(isnan(asin(M_E)));

  ok(isnan(asin(INFINITY)));
  ok(isnan(asin(-INFINITY)));
  ok(isnan(asin(NAN)));

  diag("atan");

  eq_ok(atan(0), 0);
  eq_ok(atan(1), 0.785398);
  eq_ok(atan(-1), -0.785398);
  eq_ok(atan(0.5), 0.463648);

  eq_ok(atan(1.23e300), 1.570796);
  ok(isnegzero(atan(-1.23e-300)));

  eq_ok(atan(M_PI), 1.262627);
  eq_ok(atan(M_E), 1.218283);

  eq_ok(atan(INFINITY), 1.570796);
  eq_ok(atan(-INFINITY), -1.570796);
  nan_ok(atan(NAN));

  diag("atan2");

  // x, 0
  eq_ok(atan2(0, 0), 0);
  eq_ok(atan2(1, 0), 1.570796);
  eq_ok(atan2(-1, 0), -1.570796);
  eq_ok(atan2(0.5, 0), 1.570796);

  eq_ok(atan2(1.23e300, 0), 1.570796);
  ok(isnegzero(atan2(-1.23e-300, 0)));

  eq_ok(atan2(M_PI, 0), 1.570796);
  eq_ok(atan2(M_E, 0), 1.570796);

  eq_ok(atan2(INFINITY, 0), 1.570796);
  eq_ok(atan2(-INFINITY, 0), -1.570796);
  ok(isnan(atan2(NAN, 0)));

  // x, 1
  eq_ok(atan2(0, 1), 0);
  eq_ok(atan2(1, 1), 0.785398);
  eq_ok(atan2(-1, 1), -0.785398);
  eq_ok(atan2(0.5, 1), 0.463648);

  eq_ok(atan2(1.23e300, 1), 1.570796);
  ok(isnegzero(atan2(-1.23e-300, 1)));

  eq_ok(atan2(M_PI, 1), 1.262627);
  eq_ok(atan2(M_E, 1), 1.218283);

  eq_ok(atan2(INFINITY, 1), 1.570796);
  eq_ok(atan2(-INFINITY, 1), -1.570796);
  ok(isnan(atan2(NAN, 1)));

  // x, INFINITY
  eq_ok(atan2(0, INFINITY), 0);
  eq_ok(atan2(1, INFINITY), 0);
  ok(isnegzero(atan2(-1, INFINITY)));
  eq_ok(atan2(0.5, INFINITY), 0);

  eq_ok(atan2(1.23e300, INFINITY), 0);
  ok(isnegzero(atan2(-1.23e-300, INFINITY)));

  eq_ok(atan2(M_PI, INFINITY), 0);
  eq_ok(atan2(M_E, INFINITY), 0);

  eq_ok(atan2(INFINITY, INFINITY), 0.785398);
  eq_ok(atan2(-INFINITY, INFINITY), -0.785398);
  ok(isnan(atan2(NAN, INFINITY)));

  // x, -INFINITY
  eq_ok(atan2(0, -INFINITY), M_PI);
  eq_ok(atan2(1, -INFINITY), M_PI);
  ok(isnegzero(atan2(-1, -INFINITY)));
  eq_ok(atan2(0.5, -INFINITY), M_PI);

  eq_ok(atan2(1.23e300, -INFINITY), M_PI);
  ok(isnegzero(atan2(-1.23e-300, -INFINITY)));

  eq_ok(atan2(M_PI, -INFINITY), M_PI);
  eq_ok(atan2(M_E, -INFINITY), M_PI);

  eq_ok(atan2(INFINITY, -INFINITY), 2.356194);
  eq_ok(atan2(-INFINITY, -INFINITY), -2.356194);
  ok(isnan(atan2(NAN, -INFINITY)));

  // x, NAN
  ok(isnan(atan2(0, NAN)));
  ok(isnan(atan2(1, NAN)));
  ok(isnan(atan2(-1, NAN)));
  ok(isnan(atan2(0.5, NAN)));

  ok(isnan(atan2(1.23e300, NAN)));
  ok(isnan(atan2(-1.23e-300, NAN)));

  ok(isnan(atan2(M_PI, NAN)));
  ok(isnan(atan2(M_E, NAN)));

  ok(isnan(atan2(INFINITY, NAN)));
  ok(isnan(atan2(-INFINITY, NAN)));
  ok(isnan(atan2(NAN, NAN)));

  diag("ceil");

  eq_ok(ceil(0), 0);
  eq_ok(ceil(1), 1);
  eq_ok(ceil(-1), -1);
  eq_ok(ceil(0.5), 1);

  eq_ok(ceil(1.23e300), 1.23e300);
  eq_ok(ceil(-1.23e-300), 0);

  eq_ok(ceil(M_PI), 4);
  eq_ok(ceil(M_E), 3);

  inf_ok(ceil(INFINITY), 1);
  inf_ok(ceil(-INFINITY), -1);
  nan_ok(ceil(NAN));

  // Each of the tests are against these values:
  //
  // * Integers: 0, 1, -1
  // * Floats: 1.23e30, -1.23e-30
  // * Doubles: 1.23e300, -1.23e-300
  // * Infinities: INFINITY, -INFINITY
  // * Not a number: NAN

  diag("cos");
  eq_ok(cos(0), 1);
  eq_ok(cos(1), 0.540302);
  eq_ok(cos(-1), 0.540302);
  eq_ok(cos(1.23e30), -0.966066);
  eq_ok(cos(-1.23e-30), 1.000000);
  eq_ok(cos(1.23e300), 0.251533);
  eq_ok(cos(-1.23e-300), 1.000000);
  ok(isnan(cos(INFINITY)));
  ok(isnan(cos(-INFINITY)));
  ok(isnan(cos(NAN)));

  // test_cos();
  // test_cosh();
  // test_exp();
  // test_fabs();
  // test_floor();
  // test_fmod();
  // test_ldexp();
  // test_log();
  // test_log10();
  // test_pow();
  // test_sin();
  // test_sinh();
  // test_sqrt();
  // test_tan();
  // test_tanh();

  done_testing();
}
