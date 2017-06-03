// Test for math.h.

#include <stdio.h>
#include <math.h>
#include "tests.h"

#define PI 3.14159265
#define IS_NAN -2147483648

unsigned long long ullmax = 18446744073709551615ull;

int main()
{
  plan(359);

  // Note: There are some tests that must be disabled because they return
  // different values under different compilers. See the comment surrounding the
  // disabled() tests for more information.

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
  is_negzero(asin(-1.23e-300));
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
  is_negzero(atan(-1.23e-300));
  is_eq(atan(M_PI), 1.262627);
  is_eq(atan(M_E), 1.218283);
  is_eq(atan(INFINITY), 1.570796);
  is_eq(atan(-INFINITY), -1.570796);
  is_nan(atan(NAN));

  diag("atan2");

  // atan2(x, 0)
  is_eq(atan2(0, 0), 0);
  is_eq(atan2(1, 0), 1.570796);
  is_eq(atan2(-1, 0), -1.570796);
  is_eq(atan2(0.5, 0), 1.570796);
  is_eq(atan2(1.23e300, 0), 1.570796);
  is_negzero(atan2(-1.23e-300, 0));
  is_eq(atan2(M_PI, 0), 1.570796);
  is_eq(atan2(M_E, 0), 1.570796);
  is_eq(atan2(INFINITY, 0), 1.570796);
  is_eq(atan2(-INFINITY, 0), -1.570796);
  is_nan(atan2(NAN, 0));

  // atan2(x, 1)
  is_eq(atan2(0, 1), 0);
  is_eq(atan2(1, 1), 0.785398);
  is_eq(atan2(-1, 1), -0.785398);
  is_eq(atan2(0.5, 1), 0.463648);
  is_eq(atan2(1.23e300, 1), 1.570796);
  is_negzero(atan2(-1.23e-300, 1));
  is_eq(atan2(M_PI, 1), 1.262627);
  is_eq(atan2(M_E, 1), 1.218283);
  is_eq(atan2(INFINITY, 1), 1.570796);
  is_eq(atan2(-INFINITY, 1), -1.570796);
  is_nan(atan2(NAN, 1));

  // atan2(x, INFINITY)
  is_eq(atan2(0, INFINITY), 0);
  is_eq(atan2(1, INFINITY), 0);
  is_negzero(atan2(-1, INFINITY));
  is_eq(atan2(0.5, INFINITY), 0);
  is_eq(atan2(1.23e300, INFINITY), 0);
  is_negzero(atan2(-1.23e-300, INFINITY));
  is_eq(atan2(M_PI, INFINITY), 0);
  is_eq(atan2(M_E, INFINITY), 0);
  is_eq(atan2(INFINITY, INFINITY), 0.785398);
  is_eq(atan2(-INFINITY, INFINITY), -0.785398);
  is_nan(atan2(NAN, INFINITY));

  // atan2(x, -INFINITY)
  is_eq(atan2(0, -INFINITY), M_PI);
  is_eq(atan2(1, -INFINITY), M_PI);
  is_negzero(atan2(-1, -INFINITY));
  is_eq(atan2(0.5, -INFINITY), M_PI);
  is_eq(atan2(1.23e300, -INFINITY), M_PI);
  is_negzero(atan2(-1.23e-300, -INFINITY));
  is_eq(atan2(M_PI, -INFINITY), M_PI);
  is_eq(atan2(M_E, -INFINITY), M_PI);
  is_eq(atan2(INFINITY, -INFINITY), 2.356194);
  is_eq(atan2(-INFINITY, -INFINITY), -2.356194);
  is_nan(atan2(NAN, -INFINITY));

  // atan2(x, NAN)
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
  is_eq(cos(0.5), 0.877583);
  // https://github.com/golang/go/issues/20539
  disabled(is_eq(cos(1.23e300), 0.251533));
  is_eq(cos(-1.23e-300), 1);
  is_eq(cos(M_PI), -1);
  is_eq(cos(M_E), -0.911734);
  is_nan(cos(INFINITY));
  is_nan(cos(-INFINITY));
  is_nan(cos(NAN));

  diag("cosh");
  is_eq(cosh(0), 1);
  is_eq(cosh(1), 1.543081);
  is_eq(cosh(-1), 1.543081);
  is_eq(cosh(0.5), 1.127626);
  // https://github.com/golang/go/issues/20539
  disabled(is_eq(cosh(1.23e300), 1));
  is_eq(cosh(-1.23e-300), 1);
  is_eq(cosh(M_PI), 11.591953);
  is_eq(cosh(M_E), 7.610125);
  is_inf(cosh(INFINITY), 1);
  is_inf(cosh(-INFINITY), 1);
  is_nan(cosh(NAN));

  diag("exp");
  is_eq(exp(0), 1);
  is_eq(exp(1), 2.718282);
  is_eq(exp(-1), 0.367879);
  is_eq(exp(0.5), 1.648721);
  // https://github.com/golang/go/issues/20539
  disabled(is_inf(exp(1.23e300), 1));
  is_eq(exp(-1.23e-300), 1);
  is_eq(exp(M_PI), 23.140693);
  is_eq(exp(M_E), 15.154262);
  is_inf(exp(INFINITY), 1);
  is_eq(exp(-INFINITY), 0);
  is_nan(exp(NAN));

  diag("fabs");
  is_eq(fabs(0), 0);
  is_eq(fabs(1), 1);
  is_eq(fabs(-1), 1);
  is_eq(fabs(0.5), 0.5);
  is_eq(fabs(1.23e300), 1.23e300);
  is_eq(fabs(-1.23e-300), 1.23e-300);
  is_eq(fabs(M_PI), M_PI);
  is_eq(fabs(M_E), M_E);
  is_inf(fabs(INFINITY), 1);
  is_inf(fabs(-INFINITY), 1);
  is_nan(fabs(NAN));

  diag("floor");
  is_eq(floor(0), 0);
  is_eq(floor(1), 1);
  is_eq(floor(-1), -1);
  is_eq(floor(0.5), 0);
  is_eq(floor(1.23e300), 1.23e300);
  is_eq(floor(-1.23e-300), -1);
  is_eq(floor(M_PI), 3);
  is_eq(floor(M_E), 2);
  is_inf(floor(INFINITY), 1);
  is_inf(floor(-INFINITY), -1);
  is_nan(floor(NAN));

  diag("fmod");

  // fmod(x, 0)
  is_nan(fmod(0, 0));
  is_nan(fmod(1, 0));
  is_nan(fmod(-1, 0));
  is_nan(fmod(0.5, 0));
  is_nan(fmod(1.23e300, 0));
  is_nan(fmod(-1.23e-300, 0));
  is_nan(fmod(M_PI, 0));
  is_nan(fmod(M_E, 0));
  is_nan(fmod(INFINITY, 0));
  is_nan(fmod(-INFINITY, 0));
  is_nan(fmod(NAN, 0));

  // fmod(x, 0.5)
  is_eq(fmod(0, 0.5), 0);
  is_eq(fmod(1, 0.5), 0);
  is_negzero(fmod(-1, 0.5));
  is_eq(fmod(0.5, 0.5), 0);
  is_eq(fmod(1.23e300, 0.5), 0);
  is_negzero(fmod(-1.23e-300, 0.5));
  is_eq(fmod(M_PI, 0.5), M_PI - 3);
  is_eq(fmod(M_E, 0.5), M_E - 2.5);
  is_nan(fmod(INFINITY, 0.5));
  is_nan(fmod(-INFINITY, 0.5));
  is_nan(fmod(NAN, 0.5));

  // fmod(x, INFINITY)
  is_eq(fmod(0, INFINITY), 0);
  is_eq(fmod(1, INFINITY), 1);
  is_negzero(fmod(-1, INFINITY));
  is_eq(fmod(0.5, INFINITY), 0.5);
  is_eq(fmod(1.23e300, INFINITY), 1.23e300);
  is_negzero(fmod(-1.23e-300, INFINITY));
  is_eq(fmod(M_PI, INFINITY), M_PI);
  is_eq(fmod(M_E, INFINITY), M_E);
  is_nan(fmod(INFINITY, INFINITY));
  is_nan(fmod(-INFINITY, INFINITY));
  is_nan(fmod(NAN, INFINITY));

  // fmod(x, -INFINITY)
  is_eq(fmod(0, -INFINITY), 0);
  is_eq(fmod(1, -INFINITY), 1);
  is_negzero(fmod(-1, -INFINITY));
  is_eq(fmod(0.5, -INFINITY), 0.5);
  is_eq(fmod(1.23e300, -INFINITY), 1.23e300);
  is_negzero(fmod(-1.23e-300, -INFINITY));
  is_eq(fmod(M_PI, -INFINITY), M_PI);
  is_eq(fmod(M_E, -INFINITY), M_E);
  is_nan(fmod(INFINITY, -INFINITY));
  is_nan(fmod(-INFINITY, -INFINITY));
  is_nan(fmod(NAN, -INFINITY));

  // fmod(x, NAN)
  is_nan(fmod(0, NAN));
  is_nan(fmod(1, NAN));
  is_nan(fmod(-1, NAN));
  is_nan(fmod(0.5, NAN));
  is_nan(fmod(1.23e300, NAN));
  is_nan(fmod(-1.23e-300, NAN));
  is_nan(fmod(M_PI, NAN));
  is_nan(fmod(M_E, NAN));
  is_nan(fmod(INFINITY, NAN));
  is_nan(fmod(-INFINITY, NAN));
  is_nan(fmod(NAN, NAN));

  diag("ldexp");
  is_eq(ldexp(0, 2), 0);
  is_eq(ldexp(1, 2), 4);
  is_eq(ldexp(-1, 2), -4);
  is_eq(ldexp(0.5, 2), 2);
  is_eq(ldexp(1.23e300, 2), 1.23e300);
  is_negzero(ldexp(-1.23e-300, 2));
  is_eq(ldexp(M_PI, 2), 12.566371);
  is_eq(ldexp(M_E, 2), 10.873127);
  is_inf(ldexp(INFINITY, 2), 1);
  is_inf(ldexp(-INFINITY, 2), -1);
  is_nan(ldexp(NAN, 2));

  diag("log");
  is_inf(log(0), -1);
  is_eq(log(1), 0);
  is_nan(log(-1));
  is_eq(log(0.5), -0.693147);
  is_eq(log(1.23e300), 1.23e300);
  is_nan(log(-1.23e-300));
  is_eq(log(M_PI), 1.144730);
  is_eq(log(M_E), 1);
  is_inf(log(INFINITY), 1);
  is_nan(log(-INFINITY));
  is_nan(log(NAN));

  diag("log10");
  is_inf(log10(0), -1);
  is_eq(log10(1), 0);
  is_nan(log10(-1));
  is_eq(log10(0.5), -0.301030);
  is_eq(log10(1.23e300), 1.23e300);
  is_nan(log10(-1.23e-300));
  is_eq(log10(M_PI), 0.497150);
  is_eq(log10(M_E), 0.434294);
  is_inf(log10(INFINITY), 1);
  is_nan(log10(-INFINITY));
  is_nan(log10(NAN));

  diag("pow");

  // pow(x, 0)
  is_eq(pow(0, 0), 1);
  is_eq(pow(1, 0), 1);
  is_eq(pow(-1, 0), 1);
  is_eq(pow(0.5, 0), 1);
  is_eq(pow(1.23e300, 0), 1);
  is_eq(pow(-1.23e-300, 0), 1);
  is_eq(pow(M_PI, 0), 1);
  is_eq(pow(M_E, 0), 1);
  is_eq(pow(INFINITY, 0), 1);
  is_eq(pow(-INFINITY, 0), 1);
  is_eq(pow(NAN, 0), 1);

  // pow(x, M_PI)
  is_eq(pow(0, M_PI), 0);
  is_eq(pow(1, M_PI), 1);
  is_nan(pow(-1, M_PI));
  is_eq(pow(0.5, M_PI), 0.113315);
  is_inf(pow(1.23e300, M_PI), 1);
  is_nan(pow(-1.23e-300, M_PI));
  is_eq(pow(M_PI, M_PI), 36.462160);
  is_eq(pow(M_E, M_PI), 23.140693);
  is_inf(pow(INFINITY, M_PI), 1);
  is_inf(pow(-INFINITY, M_PI), 1);
  is_nan(pow(NAN, M_PI));

  // pow(x, INFINITY)
  is_eq(pow(0, INFINITY), 0);
  is_eq(pow(1, INFINITY), 1);
  is_eq(pow(-1, INFINITY), 1);
  is_eq(pow(0.5, INFINITY), 0);
  is_inf(pow(1.23e300, INFINITY), 1);
  is_eq(pow(-1.23e-300, INFINITY), 0);
  is_inf(pow(M_PI, INFINITY), 1);
  is_inf(pow(M_E, INFINITY), 1);
  is_inf(pow(INFINITY, INFINITY), 1);
  is_inf(pow(-INFINITY, INFINITY), 1);
  is_nan(pow(NAN, INFINITY));

  // pow(x, -INFINITY)
  is_inf(pow(0, -INFINITY), 1);
  is_eq(pow(1, -INFINITY), 1);
  is_eq(pow(-1, -INFINITY), 1);
  is_inf(pow(0.5, -INFINITY), 1);
  is_eq(pow(1.23e300, -INFINITY), 0);
  is_inf(pow(-1.23e-300, -INFINITY), 1);
  is_eq(pow(M_PI, -INFINITY), 0);
  is_eq(pow(M_E, -INFINITY), 0);
  is_eq(pow(INFINITY, -INFINITY), 0);
  is_eq(pow(-INFINITY, -INFINITY), 0);
  is_nan(pow(NAN, -INFINITY));

  // pow(x, NAN)
  is_nan(pow(0, NAN));
  is_eq(pow(1, NAN), 1);
  is_nan(pow(-1, NAN));
  is_nan(pow(0.5, NAN));
  is_nan(pow(1.23e300, NAN));
  is_nan(pow(-1.23e-300, NAN));
  is_nan(pow(M_PI, NAN));
  is_nan(pow(M_E, NAN));
  is_nan(pow(INFINITY, NAN));
  is_nan(pow(-INFINITY, NAN));
  is_nan(pow(NAN, NAN));

  diag("sin");
  is_eq(sin(0), 0);
  is_eq(sin(1), 0.841471);
  is_eq(sin(-1), -0.841471);
  is_eq(sin(0.5), 0.479426);
  // https://github.com/golang/go/issues/20539
  disabled(is_eq(sin(1.23e300), 0.967849));
  is_negzero(sin(-1.23e-300));
  is_eq(sin(M_PI), 0);
  is_eq(sin(M_E), 0.410781);
  is_nan(sin(INFINITY));
  is_nan(sin(-INFINITY));
  is_nan(sin(NAN));

  diag("sinh");
  is_eq(sinh(0), 0);
  is_eq(sinh(1), 1.175201);
  is_eq(sinh(-1), -1.175201);
  is_eq(sinh(0.5), 0.521095);
  // https://github.com/golang/go/issues/20539
  disabled(is_eq(sinh(1.23e300), 1));
  is_negzero(sinh(-1.23e-300));
  is_eq(sinh(M_PI), 11.548739);
  is_eq(sinh(M_E), 7.544137);
  is_inf(sinh(INFINITY), 1);
  is_inf(sinh(-INFINITY), -1);
  is_nan(sinh(NAN));

  diag("sqrt");
  is_eq(sqrt(0), 0);
  is_eq(sqrt(1), 1);
  is_nan(sqrt(-1));
  is_eq(sqrt(0.5), 0.707107);
  is_eq(sqrt(1.23e300), 1.109054e150);
  is_nan(sqrt(-1.23e-300));
  is_eq(sqrt(M_PI), 1.772454);
  is_eq(sqrt(M_E), 1.648721);
  is_inf(sqrt(INFINITY), 1);
  is_nan(sqrt(-INFINITY));
  is_nan(sqrt(NAN));

  diag("tan");
  is_eq(tan(0), 0);
  is_eq(tan(1), 1.557408);
  is_eq(tan(-1), -1.557408);
  is_eq(tan(0.5), 0.546302);
  // https://github.com/golang/go/issues/20539
  disabled(is_eq(tan(1.23e300), 3.847798));
  is_negzero(tan(-1.23e-300));
  is_eq(tan(M_PI), 0);
  is_eq(tan(M_E), -0.450550);
  is_nan(tan(INFINITY));
  is_nan(tan(-INFINITY));
  is_nan(tan(NAN));

  diag("tanh");
  is_eq(tanh(0), 0);
  is_eq(tanh(1), 0.761594);
  is_eq(tanh(-1), -0.761594);
  is_eq(tanh(0.5), 0.462117);
  is_eq(tanh(1.23e300), 1);
  is_negzero(tanh(-1.23e-300));
  is_eq(tanh(M_PI), 0.996272);
  is_eq(tanh(M_E), 0.991329);
  is_eq(tanh(INFINITY), 1);
  is_eq(tanh(-INFINITY), -1);
  is_nan(tanh(NAN));

  done_testing();
}
