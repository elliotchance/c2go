#include <stdio.h>
#include "tests.h"

enum number{zero, one, two, three};

enum year{Jan, Feb, Mar, Apr, May, Jun, Jul, 
          Aug, Sep, Oct, Nov, Dec};

enum State {Working = 1, Failed = 0, Freezed = 0};

enum day {sunday = 1, monday, tuesday = 5,
          wednesday, thursday = 10, friday, saturday}; 

enum state {WORKING = 0, FAILED, FREEZED};
enum state currState = 2;
enum state FindState() {return currState;}

// main function

int main()
{
	plan(15);

	// step 
	enum number n;
	n = two;
	is_eq(two ,2);
	is_eq(n   ,2);

	// step 
	for (int i=Jan; i < Feb; i++){   
		is_eq(i, Jan);
	}

	// step 
	is_eq( Working , 1);
	is_eq( Failed  , 0);
	is_eq( Freezed , 0);

	// step 
	enum day d = thursday;
	is_eq( d , 10);

	// step 
	is_eq( sunday    ,  1);
	is_eq( monday    ,  2);
	is_eq( tuesday   ,  5);
	is_eq( wednesday ,  6);
	is_eq( thursday  , 10);
	is_eq( friday    , 11);
	is_eq( saturday  , 12);

	// step 
	is_eq( FindState() , FREEZED);

	done_testing();
}
