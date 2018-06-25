// +build windows nacl plan9
// A nop syslog logger for platforms with no syslog support

package noarch

// void    closelog(void);
func Closelog() {
}

// void    openlog(const char *, int, int);
func Openlog(ident *byte, logopt int32, facility int32) {
}

// int     setlogmask(int);
func Setlogmask(mask int32) int32 {
	return 0
}

// void    syslog(int, const char *, ...);
func Syslog(priority int32, format *byte, args ...interface{}) {
}

// void    vsyslog(int, const char *, struct __va_list_tag *);
func Vsyslog(priority int32, format *byte, args VaList) {
}
