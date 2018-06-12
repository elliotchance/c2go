#include <stdio.h>
#include "tests.h"

/* Text */
enum number{zero, one, two, three};

/**
 * Text
 */
enum
{
  _ISupper = ((0) < 8 ? ((1 << (0)) << 8) : ((1 << (0)) >> 8)),
  _ISalnum = ((11) < 8 ? ((1 << (11)) << 8) : ((1 << (11)) >> 8))
};

/** <b> Text </b> */
enum year{Jan, Feb, Mar, Apr, May, Jun, Jul, 
          Aug, Sep, Oct, Nov, Dec};

// Text
enum State {Working = 1, Failed = 0, Freezed = 0};

// Text
// Text
enum day {sunday = 1, monday, tuesday = 5,
          wednesday, thursday = 10, friday, saturday}; 

enum state {WORKING = 0, FAILED, FREEZED};
enum state currState = 2;
enum state FindState() {return currState;}

enum { FLY , JUMP };

/// TYPEDEF
typedef enum {
    a, b, c
 } T_ENUM;

/**
 * Text 
 */
typedef enum e_strategy {RANDOM, IMMEDIATE = 5, SEARCH} strategy;

enum { ESC_A = 1, ESC_d };

// main function

int main()
{
	plan(30);

	// step 1
	enum number n;
	n = two;
	is_eq(two ,2);
	is_eq(n   ,2);

	// step 2
	is_eq(_ISupper ,256);
	is_eq(_ISalnum ,8  );

	// step 3
	for (int i=Jan; i < Feb; i++){   
		is_eq(i, Jan);
	}

	// step 4
	is_eq( Working , 1);
	is_eq( Failed  , 0);
	is_eq( Freezed , 0);

	// step 5
	enum day d = thursday;
	is_eq( d , 10);

	// step 6
	is_eq( sunday    ,  1);
	is_eq( monday    ,  2);
	is_eq( tuesday   ,  5);
	is_eq( wednesday ,  6);
	is_eq( thursday  , 10);
	is_eq( friday    , 11);
	is_eq( saturday  , 12);

	// step 7
	is_eq( FindState() , FREEZED);

	// step
	T_ENUM cc = a;
	is_eq( cc , a );
	cc = c;
	is_eq( cc , c );
	cc = (T_ENUM)(1);
	is_eq( cc , b );

	// step
	strategy str = RANDOM;
	is_eq( str , RANDOM );
	enum e_strategy e_str = RANDOM;
	is_eq( e_str, RANDOM );
	is_eq( str , e_str );
	is_eq(IMMEDIATE , 5);
	is_eq(SEARCH    , 6);

	// step 
	is_eq( FLY  , 0 );
	is_eq( JUMP , 1 );

	is_eq(ESC_d, 2);

	diag("sizeof")
	is_eq(sizeof(JUMP ),sizeof(int));
	is_eq(sizeof(Jan  ),sizeof(int));

	done_testing();
}
