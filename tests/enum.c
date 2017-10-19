#include <stdio.h>
#include "tests.h"

enum number{zero, one, two, three};

enum
{
  _ISupper = ((0) < 8 ? ((1 << (0)) << 8) : ((1 << (0)) >> 8)),
  _ISalnum = ((11) < 8 ? ((1 << (11)) << 8) : ((1 << (11)) >> 8))
};

int main()
{
	plan(4);
    
	enum number n;
	n = two;
	is_eq(two ,2);
	is_eq(n   ,2);
	
	is_eq(_ISupper ,256);
	is_eq(_ISalnum ,8  );

	done_testing();
}
