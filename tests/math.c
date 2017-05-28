// Test for math.h.

#include <stdio.h>
#include <math.h>
#include "tests.h"

#define PI 3.14159265
#define IS_NAN -2147483648

unsigned long long ullmax = 18446744073709551615ull;

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
  is_eq(M_E, 2.718282);
  is_eq(M_LOG2E, 1.442695);
  is_eq(M_LOG10E, 0.434294);
  is_eq(M_LN2, 0.693147);
  is_eq(M_LN10, 2.302585);
  is_eq(M_PI, 3.141593);
  is_eq(M_PI_2, 1.570796);
  is_eq(M_PI_4, 0.785398);
  is_eq(M_1_PI, 0.318310);
  is_eq(M_2_PI, 0.636620);
  is_eq(M_2_SQRTPI, 1.128379);
  is_eq(M_SQRT2, 1.414214);
  is_eq(M_SQRT1_2, 0.707107);

  // Each of the tests are against these values:
  //
  // * Simple: 0, 1, -1, 0.5
  // * Large and small: 1.23e300, -1.23e-300
  // * Constants: M_PI, M_E
  // * Special: INFINITY, -INFINITY, NAN

  diag("acos");

  is_eq(acos(0), 1.570796);
  is_eq(acos(1), 0);
  is_eq(acos(-1), M_PI);
  is_eq(acos(0.5), 1.047198);

  is_nan(acos(1.23e300));
  is_eq(acos(-1.23e-300), 1.570796);

  is_nan(acos(M_PI));
  is_nan(acos(M_E));

  is_nan(acos(INFINITY));
  is_nan(acos(-INFINITY));
  is_nan(acos(NAN));

  diag("asin");

  is_eq(asin(0), 0);
  is_eq(asin(1), 1.570796);
  is_eq(asin(-1), -1.570796);
  is_eq(asin(0.5), 0.523599);

  is_nan(asin(1.23e300));
  ok(isnegzero(asin(-1.23e-300)));

  is_nan(asin(M_PI));
  is_nan(asin(M_E));

  is_nan(asin(INFINITY));
  is_nan(asin(-INFINITY));
  is_nan(asin(NAN));

  diag("atan");

  is_eq(atan(0), 0);
  is_eq(atan(1), 0.785398);
  is_eq(atan(-1), -0.785398);
  is_eq(atan(0.5), 0.463648);

  is_eq(atan(1.23e300), 1.570796);
  ok(isnegzero(atan(-1.23e-300)));

  is_eq(atan(M_PI), 1.262627);
  is_eq(atan(M_E), 1.218283);

  is_eq(atan(INFINITY), 1.570796);
  is_eq(atan(-INFINITY), -1.570796);
  is_nan(atan(NAN));

  diag("atan2");

  // x, 0
  is_eq(atan2(0, 0), 0);
  is_eq(atan2(1, 0), 1.570796);
  is_eq(atan2(-1, 0), -1.570796);
  is_eq(atan2(0.5, 0), 1.570796);

  is_eq(atan2(1.23e300, 0), 1.570796);
  ok(isnegzero(atan2(-1.23e-300, 0)));

  is_eq(atan2(M_PI, 0), 1.570796);
  is_eq(atan2(M_E, 0), 1.570796);

  is_eq(atan2(INFINITY, 0), 1.570796);
  is_eq(atan2(-INFINITY, 0), -1.570796);
  is_nan(atan2(NAN, 0));

  // x, 1
  is_eq(atan2(0, 1), 0);
  is_eq(atan2(1, 1), 0.785398);
  is_eq(atan2(-1, 1), -0.785398);
  is_eq(atan2(0.5, 1), 0.463648);

  is_eq(atan2(1.23e300, 1), 1.570796);
  ok(isnegzero(atan2(-1.23e-300, 1)));

  is_eq(atan2(M_PI, 1), 1.262627);
  is_eq(atan2(M_E, 1), 1.218283);

  is_eq(atan2(INFINITY, 1), 1.570796);
  is_eq(atan2(-INFINITY, 1), -1.570796);
  is_nan(atan2(NAN, 1));

  // x, INFINITY
  is_eq(atan2(0, INFINITY), 0);
  is_eq(atan2(1, INFINITY), 0);
  ok(isnegzero(atan2(-1, INFINITY)));
  is_eq(atan2(0.5, INFINITY), 0);

  is_eq(atan2(1.23e300, INFINITY), 0);
  ok(isnegzero(atan2(-1.23e-300, INFINITY)));

  is_eq(atan2(M_PI, INFINITY), 0);
  is_eq(atan2(M_E, INFINITY), 0);

  is_eq(atan2(INFINITY, INFINITY), 0.785398);
  is_eq(atan2(-INFINITY, INFINITY), -0.785398);
  is_nan(atan2(NAN, INFINITY));

  // x, -INFINITY
  is_eq(atan2(0, -INFINITY), M_PI);
  is_eq(atan2(1, -INFINITY), M_PI);
  ok(isnegzero(atan2(-1, -INFINITY)));
  is_eq(atan2(0.5, -INFINITY), M_PI);

  is_eq(atan2(1.23e300, -INFINITY), M_PI);
  ok(isnegzero(atan2(-1.23e-300, -INFINITY)));

  is_eq(atan2(M_PI, -INFINITY), M_PI);
  is_eq(atan2(M_E, -INFINITY), M_PI);

  is_eq(atan2(INFINITY, -INFINITY), 2.356194);
  is_eq(atan2(-INFINITY, -INFINITY), -2.356194);
  is_nan(atan2(NAN, -INFINITY));

  // x, NAN
  is_nan(atan2(0, NAN));
  is_nan(atan2(1, NAN));
  is_nan(atan2(-1, NAN));
  is_nan(atan2(0.5, NAN));

  is_nan(atan2(1.23e300, NAN));
  is_nan(atan2(-1.23e-300, NAN));

  is_nan(atan2(M_PI, NAN));
  is_nan(atan2(M_E, NAN));

  is_nan(atan2(INFINITY, NAN));
  is_nan(atan2(-INFINITY, NAN));
  is_nan(atan2(NAN, NAN));

  diag("ceil");

  is_eq(ceil(0), 0);
  is_eq(ceil(1), 1);
  is_eq(ceil(-1), -1);
  is_eq(ceil(0.5), 1);

  is_eq(ceil(1.23e300), 1.23e300);
  is_eq(ceil(-1.23e-300), 0);

  is_eq(ceil(M_PI), 4);
  is_eq(ceil(M_E), 3);

  is_inf(ceil(INFINITY), 1);
  is_inf(ceil(-INFINITY), -1);
  is_nan(ceil(NAN));

  diag("cos");
  is_eq(cos(0), 1);
  is_eq(cos(1), 0.540302);
  is_eq(cos(-1), 0.540302);
  is_eq(cos(1.23e30), -0.966066);
  is_eq(cos(-1.23e-30), 1.000000);
  is_eq(cos(1.23e300), 0.251533);
  is_eq(cos(-1.23e-300), 1.000000);
  is_nan(cos(INFINITY));
  is_nan(cos(-INFINITY));
  is_nan(cos(NAN));

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
