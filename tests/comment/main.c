// comment1

/* comment2 */

/*
 * comment3
 */

void /* comment4 */ a /* comment5 */( /* comment6 */ int /* comment7 */ i /*comment8*/) // comment9
{
	// comment10
	// comment11
	(/* comment12 */ void /* comment13 */)/*comment14*/(/*comment15*/ i /*comment16 */)/*comment17*/;
}
// comment18
//comment19

void b //comment20
( // comment21
 )//comment22
{ //comment23
//comment24
}//comment25


void /* comment26 */ main /*comment27*/()
{
	int i = 0;
	for ( i = 0 ; i < 5 ; i++)
	{
		if (i > 2)
		{
			a(i);
		} else {
			/*
			 * * * // comment28
			*/
			b();
		}
	}
}
/*
 * comment29
 *
 *
 * 
 */
// comment30

