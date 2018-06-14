package noarch

import (
	"strings"
)

const (
	EPERM   = 1  // Operation not permitted
	ENOENT  = 2  // No such file or directory
	ESRCH   = 3  // No such process
	EINTR   = 4  // Interrupted system call
	EIO     = 5  // I/O error
	ENXIO   = 6  // No such device or address
	E2BIG   = 7  // Argument list too long
	ENOEXEC = 8  // Exec format error
	EBADF   = 9  // Bad file number
	ECHILD  = 10 // No child processes
	EAGAIN  = 11 // Try again
	ENOMEM  = 12 // Out of memory
	EACCES  = 13 // Permission denied
	EFAULT  = 14 // Bad address
	ENOTBLK = 15 // Block device required
	EBUSY   = 16 // Device or resource busy
	EEXIST  = 17 // File exists
	EXDEV   = 18 // Cross-device link
	ENODEV  = 19 // No such device
	ENOTDIR = 20 // Not a directory
	EISDIR  = 21 // Is a directory
	EINVAL  = 22 // Invalid argument
	ENFILE  = 23 // File table overflow
	EMFILE  = 24 // Too many open files
	ENOTTY  = 25 // Not a typewriter
	ETXTBSY = 26 // Text file busy
	EFBIG   = 27 // File too large
	ENOSPC  = 28 // No space left on device
	ESPIPE  = 29 // Illegal seek
	EROFS   = 30 // Read-only file system
	EMLINK  = 31 // Too many links
	EPIPE   = 32 // Broken pipe
	EDOM    = 33 // Math argument out of domain of func
	ERANGE  = 34 // Math result not representable

	EDEADLK      = 35     // Resource deadlock would occur
	ENAMETOOLONG = 36     // File name too long
	ENOLCK       = 37     // No record locks available
	ENOSYS       = 38     // Function not implemented
	ENOTEMPTY    = 39     // Directory not empty
	ELOOP        = 40     // Too many symbolic links encountered
	EWOULDBLOCK  = EAGAIN // Operation would block
	ENOMSG       = 42     // No message of desired type
	EIDRM        = 43     // Identifier removed
	ECHRNG       = 44     // Channel number out of range
	EL2NSYNC     = 45     // Level 2 not synchronized
	EL3HLT       = 46     // Level 3 halted
	EL3RST       = 47     // Level 3 reset
	ELNRNG       = 48     // Link number out of range
	EUNATCH      = 49     // Protocol driver not attached
	ENOCSI       = 50     // No CSI structure available
	EL2HLT       = 51     // Level 2 halted
	EBADE        = 52     // Invalid exchange
	EBADR        = 53     // Invalid request descriptor
	EXFULL       = 54     // Exchange full
	ENOANO       = 55     // No anode
	EBADRQC      = 56     // Invalid request code
	EBADSLT      = 57     // Invalid slot

	EDEADLOCK = EDEADLK

	EBFONT          = 59  // Bad font file format
	ENOSTR          = 60  // Device not a stream
	ENODATA         = 61  // No data available
	ETIME           = 62  // Timer expired
	ENOSR           = 63  // Out of streams resources
	ENONET          = 64  // Machine is not on the network
	ENOPKG          = 65  // Package not installed
	EREMOTE         = 66  // Object is remote
	ENOLINK         = 67  // Link has been severed
	EADV            = 68  // Advertise error
	ESRMNT          = 69  // Srmount error
	ECOMM           = 70  // Communication error on send
	EPROTO          = 71  // Protocol error
	EMULTIHOP       = 72  // Multihop attempted
	EDOTDOT         = 73  // RFS specific error
	EBADMSG         = 74  // Not a data message
	EOVERFLOW       = 75  // Value too large for defined data type
	ENOTUNIQ        = 76  // Name not unique on network
	EBADFD          = 77  // File descriptor in bad state
	EREMCHG         = 78  // Remote address changed
	ELIBACC         = 79  // Can not access a needed shared library
	ELIBBAD         = 80  // Accessing a corrupted shared library
	ELIBSCN         = 81  // .lib section in a.out corrupted
	ELIBMAX         = 82  // Attempting to link in too many shared libraries
	ELIBEXEC        = 83  // Cannot exec a shared library directly
	EILSEQ          = 84  // Illegal byte sequence
	ERESTART        = 85  // Interrupted system call should be restarted
	ESTRPIPE        = 86  // Streams pipe error
	EUSERS          = 87  // Too many users
	ENOTSOCK        = 88  // Socket operation on non-socket
	EDESTADDRREQ    = 89  // Destination address required
	EMSGSIZE        = 90  // Message too long
	EPROTOTYPE      = 91  // Protocol wrong type for socket
	ENOPROTOOPT     = 92  // Protocol not available
	EPROTONOSUPPORT = 93  // Protocol not supported
	ESOCKTNOSUPPORT = 94  // Socket type not supported
	EOPNOTSUPP      = 95  // Operation not supported on transport endpoint
	EPFNOSUPPORT    = 96  // Protocol family not supported
	EAFNOSUPPORT    = 97  // Address family not supported by protocol
	EADDRINUSE      = 98  // Address already in use
	EADDRNOTAVAIL   = 99  // Cannot assign requested address
	ENETDOWN        = 100 // Network is down
	ENETUNREACH     = 101 // Network is unreachable
	ENETRESET       = 102 // Network dropped connection because of reset
	ECONNABORTED    = 103 // Software caused connection abort
	ECONNRESET      = 104 // Connection reset by peer
	ENOBUFS         = 105 // No buffer space available
	EISCONN         = 106 // Transport endpoint is already connected
	ENOTCONN        = 107 // Transport endpoint is not connected
	ESHUTDOWN       = 108 // Cannot send after transport endpoint shutdown
	ETOOMANYREFS    = 109 // Too many references: cannot splice
	ETIMEDOUT       = 110 // Connection timed out
	ECONNREFUSED    = 111 // Connection refused
	EHOSTDOWN       = 112 // Host is down
	EHOSTUNREACH    = 113 // No route to host
	EALREADY        = 114 // Operation already in progress
	EINPROGRESS     = 115 // Operation now in progress
	ESTALE          = 116 // Stale NFS file handle
	EUCLEAN         = 117 // Structure needs cleaning
	ENOTNAM         = 118 // Not a XENIX named type file
	ENAVAIL         = 119 // No XENIX semaphores available
	EISNAM          = 120 // Is a named type file
	EREMOTEIO       = 121 // Remote I/O error
	EDQUOT          = 122 // Quota exceeded

	ENOMEDIUM    = 123 // No medium found
	EMEDIUMTYPE  = 124 // Wrong medium type
	ECANCELED    = 125 // Operation Canceled
	ENOKEY       = 126 // Required key not available
	EKEYEXPIRED  = 127 // Key has expired
	EKEYREVOKED  = 128 // Key has been revoked
	EKEYREJECTED = 129 // Key was rejected by service

	// for robust mutexes
	EOWNERDEAD      = 130 // Owner died
	ENOTRECOVERABLE = 131 // State not recoverable
)

// Error table
var errors = [...]string{
	1:   "operation not permitted",
	2:   "no such file or directory",
	3:   "no such process",
	4:   "interrupted system call",
	5:   "input/output error",
	6:   "no such device or address",
	7:   "argument list too long",
	8:   "exec format error",
	9:   "bad file descriptor",
	10:  "no child processes",
	11:  "resource temporarily unavailable",
	12:  "cannot allocate memory",
	13:  "permission denied",
	14:  "bad address",
	15:  "block device required",
	16:  "device or resource busy",
	17:  "file exists",
	18:  "invalid cross-device link",
	19:  "no such device",
	20:  "not a directory",
	21:  "is a directory",
	22:  "invalid argument",
	23:  "too many open files in system",
	24:  "too many open files",
	25:  "inappropriate ioctl for device",
	26:  "text file busy",
	27:  "file too large",
	28:  "no space left on device",
	29:  "illegal seek",
	30:  "read-only file system",
	31:  "too many links",
	32:  "broken pipe",
	33:  "numerical argument out of domain",
	34:  "numerical result out of range",
	35:  "resource deadlock avoided",
	36:  "file name too long",
	37:  "no locks available",
	38:  "function not implemented",
	39:  "directory not empty",
	40:  "too many levels of symbolic links",
	42:  "no message of desired type",
	43:  "identifier removed",
	44:  "channel number out of range",
	45:  "level 2 not synchronized",
	46:  "level 3 halted",
	47:  "level 3 reset",
	48:  "link number out of range",
	49:  "protocol driver not attached",
	50:  "no CSI structure available",
	51:  "level 2 halted",
	52:  "invalid exchange",
	53:  "invalid request descriptor",
	54:  "exchange full",
	55:  "no anode",
	56:  "invalid request code",
	57:  "invalid slot",
	59:  "bad font file format",
	60:  "device not a stream",
	61:  "no data available",
	62:  "timer expired",
	63:  "out of streams resources",
	64:  "machine is not on the network",
	65:  "package not installed",
	66:  "object is remote",
	67:  "link has been severed",
	68:  "advertise error",
	69:  "srmount error",
	70:  "communication error on send",
	71:  "protocol error",
	72:  "multihop attempted",
	73:  "RFS specific error",
	74:  "bad message",
	75:  "value too large for defined data type",
	76:  "name not unique on network",
	77:  "file descriptor in bad state",
	78:  "remote address changed",
	79:  "can not access a needed shared library",
	80:  "accessing a corrupted shared library",
	81:  ".lib section in a.out corrupted",
	82:  "attempting to link in too many shared libraries",
	83:  "cannot exec a shared library directly",
	84:  "invalid or incomplete multibyte or wide character",
	85:  "interrupted system call should be restarted",
	86:  "streams pipe error",
	87:  "too many users",
	88:  "socket operation on non-socket",
	89:  "destination address required",
	90:  "message too long",
	91:  "protocol wrong type for socket",
	92:  "protocol not available",
	93:  "protocol not supported",
	94:  "socket type not supported",
	95:  "operation not supported",
	96:  "protocol family not supported",
	97:  "address family not supported by protocol",
	98:  "address already in use",
	99:  "cannot assign requested address",
	100: "network is down",
	101: "network is unreachable",
	102: "network dropped connection on reset",
	103: "software caused connection abort",
	104: "connection reset by peer",
	105: "no buffer space available",
	106: "transport endpoint is already connected",
	107: "transport endpoint is not connected",
	108: "cannot send after transport endpoint shutdown",
	109: "too many references: cannot splice",
	110: "connection timed out",
	111: "connection refused",
	112: "host is down",
	113: "no route to host",
	114: "operation already in progress",
	115: "operation now in progress",
	116: "stale NFS file handle",
	117: "structure needs cleaning",
	118: "not a XENIX named type file",
	119: "no XENIX semaphores available",
	120: "is a named type file",
	121: "remote I/O error",
	122: "disk quota exceeded",
	123: "no medium found",
	124: "wrong medium type",
	125: "operation canceled",
	126: "required key not available",
	127: "key has expired",
	128: "key has been revoked",
	129: "key was rejected by service",
	130: "owner died",
	131: "state not recoverable",
	132: "operation not possible due to RF-kill",
}

var err2bytes = make(map[int][]byte)
var errInverse = make(map[string]int)

func init() {
	err2bytes[0] = []byte{0}
	for i, str := range errors {
		err2bytes[i] = []byte(str)
		err2bytes[i] = append(err2bytes[i], 0)
		err2bytes[i][0] = strings.ToUpper(string(err2bytes[i][0:1]))[0]
		errInverse[strings.ToLower(str)] = i
	}
}

// Strerror translates an errno error code into an error message.
func Strerror(errno int32) *byte {
	b, ok := err2bytes[int(errno)]
	if ok {
		return &b[0]
	}
	return &err2bytes[0][0]
}

var currentErrno int32

func setCurrentErrnoErr(err error) {
	if err == nil {
		currentErrno = 0
	} else {
		errStr := strings.ToLower(err.Error())
		if i, ok := errInverse[errStr]; ok {
			currentErrno = int32(i)
		} else {
			currentErrno = ENODATA
		}
	}
}

func setCurrentErrno(errno int32) {
	currentErrno = errno
}

// Errno returns a pointer to the current errno.
func Errno() *int32 {
	return &currentErrno
}
