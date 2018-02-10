void if_1()
{
	int a = 5;
	int b = 2;
	int c = 4;
	if ( a > b )
	{
		return;
	}
	else if ( c <= a) 
	{
		a = 0;
	}
	(void)(a);
	(void)(b);
	(void)(c);
}
