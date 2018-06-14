// +build !windows,!nacl,!plan9

package noarch

import (
	"fmt"
	"log/syslog"
)

// structure to hold information about an open logger.
var logger struct {
	ident    string
	logopt   int32
	facility syslog.Priority
	mask     int32
	w        *syslog.Writer
}

// void    closelog(void);
func Closelog() {
	logger.w.Close()
}

// void    openlog(const char *, int, int);
// TODO: handle option parameter
func Openlog(ident *byte, logopt int32, facility int32) {
	logger.ident = CStringToString(ident)
	logger.logopt = logopt // not sure what to do with this yet if anything
	logger.facility = syslog.Priority(facility)
	logger.w, _ = syslog.New(logger.facility, logger.ident)
}

// int     setlogmask(int);
// TODO
func Setlogmask(mask int32) int32 {
	ret := logger.mask
	logger.mask = mask
	return ret
}

// void    syslog(int, const char *, ...);
func Syslog(priority int32, format *byte, args ...interface{}) {
	realArgs := []interface{}{}
	realArgs = append(realArgs, convert(args)...)
	msg := fmt.Sprintf(CStringToString(format), realArgs...)
	internalSyslog(priority, msg)
}

// void    vsyslog(int, const char *, struct __va_list_tag *);
func Vsyslog(priority int32, format *byte, args VaList) {
	realArgs := []interface{}{}
	realArgs = append(realArgs, convert(args.Args)...)
	msg := fmt.Sprintf(CStringToString(format), realArgs...)
	internalSyslog(priority, msg)
}

func internalSyslog(priority int32, msg string) {
	// TODO: handle mask
	switch syslog.Priority(priority) & 0x7 { // get severity
	case syslog.LOG_EMERG:
		logger.w.Emerg(msg)
	case syslog.LOG_CRIT:
		logger.w.Crit(msg)
	case syslog.LOG_ERR:
		logger.w.Err(msg)
	case syslog.LOG_WARNING:
		logger.w.Warning(msg)
	case syslog.LOG_NOTICE:
		logger.w.Notice(msg)
	case syslog.LOG_INFO:
		logger.w.Info(msg)
	case syslog.LOG_DEBUG:
		logger.w.Debug(msg)
	}
}
