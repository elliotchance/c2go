void switch_function()
{
	int i = 34;
	switch (i) 
	{
		case (0):
		case (1):
			return;
		case (2):
			(void)(i);
			return;
		case 3:
		{
			int c;
			return;
		}
		case 4:
			break;
		case 5:
		{
			break;
		}
		case 6:
		{
		}
		case 7:
		{
			int d;
			break;
		}
	}
}
