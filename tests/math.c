// Test for math.h.

#include "tests.h"
#include <math.h>
#include <stdio.h>

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
    is_eq(M_E, 2.71828182845904509080);
    is_eq(M_LOG2E, 1.44269504088896338700);
    is_eq(M_LOG10E, 0.43429448190325181667);
    is_eq(M_LN2, 0.69314718055994528623);
    is_eq(M_LN10, 2.30258509299404590109);
    is_eq(M_PI, 3.14159265358979311600);
    is_eq(M_PI_2, 1.57079632679489655800);
    is_eq(M_PI_4, 0.78539816339744827900);
    is_eq(M_1_PI, 0.31830988618379069122);
    is_eq(M_2_PI, 0.63661977236758138243);
    is_eq(M_2_SQRTPI, 1.12837916709551255856);
    is_eq(M_SQRT2, 1.41421356237309514547);
    is_eq(M_SQRT1_2, 0.70710678118654757274);

    // Each of the tests are against these values:
    //
    // * Simple: 0, 1, -1, 0.5
    // * Large and small: 1.23e300, -1.23e-300
    // * Constants: M_PI, M_E
    // * Special: INFINITY, -INFINITY, NAN

    diag("acos");
    is_eq(acos(0), 1.57079632679489655800);
    is_eq(acos(1), 0);
    is_eq(acos(-1), M_PI);
    is_eq(acos(0.5), 1.04719755119659763132);
    is_nan(acos(1.23e300));
    is_eq(acos(-1.23e-300), 1.57079632679489655800);
    is_nan(acos(M_PI));
    is_nan(acos(M_E));
    is_nan(acos(INFINITY));
    is_nan(acos(-INFINITY));
    is_nan(acos(NAN));

    diag("asin");
    is_eq(asin(0), 0);
    is_eq(asin(1), 1.57079632679489655800);
    is_eq(asin(-1), -1.57079632679489655800);
    is_eq(asin(0.5), 0.52359877559829881566);
    is_nan(asin(1.23e300));
    is_negzero(asin(-1.23e-300));
    is_nan(asin(M_PI));
    is_nan(asin(M_E));
    is_nan(asin(INFINITY));
    is_nan(asin(-INFINITY));
    is_nan(asin(NAN));

    diag("atan");
    is_eq(atan(0), 0);
    is_eq(atan(1), 0.78539816339744827900);
    is_eq(atan(-1), -0.78539816339744827900);
    is_eq(atan(0.5), 0.46364760900080614903);
    is_eq(atan(1.23e300), 1.57079632679489655800);
    is_negzero(atan(-1.23e-300));
    is_eq(atan(M_PI), 1.26262725567891154199);
    is_eq(atan(M_E), 1.21828290501727765083);
    is_eq(atan(INFINITY), 1.57079632679489655800);
    is_eq(atan(-INFINITY), -1.57079632679489655800);
    is_nan(atan(NAN));

    diag("atan2");

    // atan2(x, 0)
    is_eq(atan2(0, 0), 0);
    is_eq(atan2(1, 0), 1.57079632679489655800);
    is_eq(atan2(-1, 0), -1.57079632679489655800);
    is_eq(atan2(0.5, 0), 1.57079632679489655800);
    is_eq(atan2(1.23e300, 0), 1.57079632679489655800);
    is_negzero(atan2(-1.23e-300, 0));
    is_eq(atan2(M_PI, 0), 1.57079632679489655800);
    is_eq(atan2(M_E, 0), 1.57079632679489655800);
    is_eq(atan2(INFINITY, 0), 1.57079632679489655800);
    is_eq(atan2(-INFINITY, 0), -1.57079632679489655800);
    is_nan(atan2(NAN, 0));

    // atan2(x, 1)
    is_eq(atan2(0, 1), 0);
    is_eq(atan2(1, 1), 0.78539816339744827900);
    is_eq(atan2(-1, 1), -0.78539816339744827900);
    is_eq(atan2(0.5, 1), 0.46364760900080609352);
    is_eq(atan2(1.23e300, 1), 1.57079632679489655800);
    is_negzero(atan2(-1.23e-300, 1));
    is_eq(atan2(M_PI, 1), 1.262627255678911764);
    is_eq(atan2(M_E, 1), 1.2182829050172776508);
    is_eq(atan2(INFINITY, 1), 1.57079632679489655800);
    is_eq(atan2(-INFINITY, 1), -1.57079632679489655800);
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
    is_eq(atan2(INFINITY, INFINITY), 0.78539816339744827900);
    is_eq(atan2(-INFINITY, INFINITY), -0.78539816339744827900);
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
    is_eq(atan2(INFINITY, -INFINITY), 2.356194490192344837);
    is_eq(atan2(-INFINITY, -INFINITY), -2.356194490192344837);
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
    is_eq(cos(1), 0.54030230586813976501);
    is_eq(cos(-1), 0.54030230586813976501);
    is_eq(cos(0.5), 0.87758256189037275874);
    // https://github.com/golang/go/issues/20539
    disabled(is_eq(cos(1.23e300), 0.251533));
    is_eq(cos(-1.23e-300), 1);
    is_eq(cos(M_PI), -1);
    is_eq(cos(M_E), -0.91173391478696508283);
    is_nan(cos(INFINITY));
    is_nan(cos(-INFINITY));
    is_nan(cos(NAN));

    diag("cosh");
    is_eq(cosh(0), 1);
    is_eq(cosh(1), 1.5430806348152437124);
    is_eq(cosh(-1), 1.5430806348152437124);
    is_eq(cosh(0.5), 1.1276259652063806982);
    // https://github.com/golang/go/issues/20539
    disabled(is_eq(cosh(1.23e300), 1));
    is_eq(cosh(-1.23e-300), 1);
    is_eq(cosh(M_PI), 11.591953275521518663);
    is_eq(cosh(M_E), 7.6101251386622870143);
    is_inf(cosh(INFINITY), 1);
    is_inf(cosh(-INFINITY), 1);
    is_nan(cosh(NAN));

    diag("exp");
    is_eq(exp(0), 1);
    is_eq(exp(1), 2.7182818284590450908);
    is_eq(exp(-1), 0.36787944117144233402);
    is_eq(exp(0.5), 1.6487212707001281942);
    // https://github.com/golang/go/issues/20539
    disabled(is_inf(exp(1.23e300), 1));
    is_eq(exp(-1.23e-300), 1);
    is_eq(exp(M_PI), 23.140692632779266802);
    is_eq(exp(M_E), 15.154262241479262485);
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
    is_eq(ldexp(1.23e300, 2), 4.92e300);
    is_negzero(ldexp(-1.23e-300, 2));
    is_eq(ldexp(M_PI, 2), 12.56637061435917246399);
    is_eq(ldexp(M_E, 2), 10.87312731383618036318);
    is_inf(ldexp(INFINITY, 2), 1);
    is_inf(ldexp(-INFINITY, 2), -1);
    is_nan(ldexp(NAN, 2));

    diag("log");
    is_inf(log(0), -1);
    is_eq(log(1), 0);
    is_nan(log(-1));
    is_eq(log(0.5), -0.69314718055994528623);
    is_eq(log(1.23e300), 690.98254206759804674221);
    is_nan(log(-1.23e-300));
    is_eq(log(M_PI), 1.14472988584940016388);
    is_eq(log(M_E), 1);
    is_inf(log(INFINITY), 1);
    is_nan(log(-INFINITY));
    is_nan(log(NAN));

    diag("log10");
    is_inf(log10(0), -1);
    is_eq(log10(1), 0);
    is_nan(log10(-1));
    is_eq(log10(0.5), -0.30102999566398119802);
    is_eq(log10(1.23e300), 300.08990511143940693728);
    is_nan(log10(-1.23e-300));
    is_eq(log10(M_PI), 0.49714987269413385418);
    is_eq(log10(M_E), 0.43429448190325181667);
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
    is_eq(pow(0.5, M_PI), 0.11331473229676088110);
    is_inf(pow(1.23e300, M_PI), 1);
    is_nan(pow(-1.23e-300, M_PI));
    is_eq(pow(M_PI, M_PI), 36.46215960720790150162);
    is_eq(pow(M_E, M_PI), 23.14069263277926324918);
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
    is_eq(sin(1), 0.84147098480789650488);
    is_eq(sin(-1), -0.84147098480789650488);
    is_eq(sin(0.5), 0.47942553860420300538);
    // https://github.com/golang/go/issues/20539
    disabled(is_eq(sin(1.23e300), 0.967849));
    is_negzero(sin(-1.23e-300));
    is_eq(sin(M_PI), 0);
    is_eq(sin(M_E), 0.41078129050290879132);
    is_nan(sin(INFINITY));
    is_nan(sin(-INFINITY));
    is_nan(sin(NAN));

    diag("sinh");
    is_eq(sinh(0), 0);
    is_eq(sinh(1), 1.1752011936438013784);
    is_eq(sinh(-1), -1.1752011936438013784);
    is_eq(sinh(0.5), 0.52109530549374738495);
    // https://github.com/golang/go/issues/20539
    disabled(is_eq(sinh(1.23e300), 1));
    is_negzero(sinh(-1.23e-300));
    is_eq(sinh(M_PI), 11.548739357257746363);
    is_eq(sinh(M_E), 7.5441371028169745827);
    is_inf(sinh(INFINITY), 1);
    is_inf(sinh(-INFINITY), -1);
    is_nan(sinh(NAN));

    diag("sqrt");
    is_eq(sqrt(0), 0);
    is_eq(sqrt(1), 1);
    is_nan(sqrt(-1));
    is_eq(sqrt(0.5), 0.70710678118654757274);
    is_eq(sqrt(1.23e300), 1.1090536506409417761e150);
    is_nan(sqrt(-1.23e-300));
    is_eq(sqrt(M_PI), 1.77245385090551588192);
    is_eq(sqrt(M_E), 1.64872127070012819416);
    is_inf(sqrt(INFINITY), 1);
    is_nan(sqrt(-INFINITY));
    is_nan(sqrt(NAN));

    diag("tan");
    is_eq(tan(0), 0);
    is_eq(tan(1), 1.55740772465490207033);
    is_eq(tan(-1), -1.55740772465490207033);
    is_eq(tan(0.5), 0.54630248984379048416);
    // https://github.com/golang/go/issues/20539
    disabled(is_eq(tan(1.23e300), 3.847798));
    is_negzero(tan(-1.23e-300));
    is_eq(tan(M_PI), 0);
    is_eq(tan(M_E), -0.45054953406980763342);
    is_nan(tan(INFINITY));
    is_nan(tan(-INFINITY));
    is_nan(tan(NAN));

    diag("tanh");
    is_eq(tanh(0), 0);
    is_eq(tanh(1), 0.76159415595576485103);
    is_eq(tanh(-1), -0.76159415595576485103);
    is_eq(tanh(0.5), 0.46211715726000973659);
    is_eq(tanh(1.23e300), 1);
    is_negzero(tanh(-1.23e-300));
    is_eq(tanh(M_PI), 0.99627207622074998028);
    is_eq(tanh(M_E), 0.99132891580059978587);
    is_eq(tanh(INFINITY), 1);
    is_eq(tanh(-INFINITY), -1);
    is_nan(tanh(NAN));

    done_testing();
}
