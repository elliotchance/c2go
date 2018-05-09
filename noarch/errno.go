package noarch

import "strings"

const EPERM = 1    /* Operation not permitted */
const ENOENT = 2   /* No such file or directory */
const ESRCH = 3    /* No such process */
const EINTR = 4    /* Interrupted system call */
const EIO = 5      /* I/O error */
const ENXIO = 6    /* No such device or address */
const E2BIG = 7    /* Argument list too long */
const ENOEXEC = 8  /* Exec format error */
const EBADF = 9    /* Bad file number */
const ECHILD = 10  /* No child processes */
const EAGAIN = 11  /* Try again */
const ENOMEM = 12  /* Out of memory */
const EACCES = 13  /* Permission denied */
const EFAULT = 14  /* Bad address */
const ENOTBLK = 15 /* Block device required */
const EBUSY = 16   /* Device or resource busy */
const EEXIST = 17  /* File exists */
const EXDEV = 18   /* Cross-device link */
const ENODEV = 19  /* No such device */
const ENOTDIR = 20 /* Not a directory */
const EISDIR = 21  /* Is a directory */
const EINVAL = 22  /* Invalid argument */
const ENFILE = 23  /* File table overflow */
const EMFILE = 24  /* Too many open files */
const ENOTTY = 25  /* Not a typewriter */
const ETXTBSY = 26 /* Text file busy */
const EFBIG = 27   /* File too large */
const ENOSPC = 28  /* No space left on device */
const ESPIPE = 29  /* Illegal seek */
const EROFS = 30   /* Read-only file system */
const EMLINK = 31  /* Too many links */
const EPIPE = 32   /* Broken pipe */
const EDOM = 33    /* Math argument out of domain of func */
const ERANGE = 34  /* Math result not representable */

const EDEADLK = 35         /* Resource deadlock would occur */
const ENAMETOOLONG = 36    /* File name too long */
const ENOLCK = 37          /* No record locks available */
const ENOSYS = 38          /* Function not implemented */
const ENOTEMPTY = 39       /* Directory not empty */
const ELOOP = 40           /* Too many symbolic links encountered */
const EWOULDBLOCK = EAGAIN /* Operation would block */
const ENOMSG = 42          /* No message of desired type */
const EIDRM = 43           /* Identifier removed */
const ECHRNG = 44          /* Channel number out of range */
const EL2NSYNC = 45        /* Level 2 not synchronized */
const EL3HLT = 46          /* Level 3 halted */
const EL3RST = 47          /* Level 3 reset */
const ELNRNG = 48          /* Link number out of range */
const EUNATCH = 49         /* Protocol driver not attached */
const ENOCSI = 50          /* No CSI structure available */
const EL2HLT = 51          /* Level 2 halted */
const EBADE = 52           /* Invalid exchange */
const EBADR = 53           /* Invalid request descriptor */
const EXFULL = 54          /* Exchange full */
const ENOANO = 55          /* No anode */
const EBADRQC = 56         /* Invalid request code */
const EBADSLT = 57         /* Invalid slot */

const EDEADLOCK = EDEADLK

const EBFONT = 59          /* Bad font file format */
const ENOSTR = 60          /* Device not a stream */
const ENODATA = 61         /* No data available */
const ETIME = 62           /* Timer expired */
const ENOSR = 63           /* Out of streams resources */
const ENONET = 64          /* Machine is not on the network */
const ENOPKG = 65          /* Package not installed */
const EREMOTE = 66         /* Object is remote */
const ENOLINK = 67         /* Link has been severed */
const EADV = 68            /* Advertise error */
const ESRMNT = 69          /* Srmount error */
const ECOMM = 70           /* Communication error on send */
const EPROTO = 71          /* Protocol error */
const EMULTIHOP = 72       /* Multihop attempted */
const EDOTDOT = 73         /* RFS specific error */
const EBADMSG = 74         /* Not a data message */
const EOVERFLOW = 75       /* Value too large for defined data type */
const ENOTUNIQ = 76        /* Name not unique on network */
const EBADFD = 77          /* File descriptor in bad state */
const EREMCHG = 78         /* Remote address changed */
const ELIBACC = 79         /* Can not access a needed shared library */
const ELIBBAD = 80         /* Accessing a corrupted shared library */
const ELIBSCN = 81         /* .lib section in a.out corrupted */
const ELIBMAX = 82         /* Attempting to link in too many shared libraries */
const ELIBEXEC = 83        /* Cannot exec a shared library directly */
const EILSEQ = 84          /* Illegal byte sequence */
const ERESTART = 85        /* Interrupted system call should be restarted */
const ESTRPIPE = 86        /* Streams pipe error */
const EUSERS = 87          /* Too many users */
const ENOTSOCK = 88        /* Socket operation on non-socket */
const EDESTADDRREQ = 89    /* Destination address required */
const EMSGSIZE = 90        /* Message too long */
const EPROTOTYPE = 91      /* Protocol wrong type for socket */
const ENOPROTOOPT = 92     /* Protocol not available */
const EPROTONOSUPPORT = 93 /* Protocol not supported */
const ESOCKTNOSUPPORT = 94 /* Socket type not supported */
const EOPNOTSUPP = 95      /* Operation not supported on transport endpoint */
const EPFNOSUPPORT = 96    /* Protocol family not supported */
const EAFNOSUPPORT = 97    /* Address family not supported by protocol */
const EADDRINUSE = 98      /* Address already in use */
const EADDRNOTAVAIL = 99   /* Cannot assign requested address */
const ENETDOWN = 100       /* Network is down */
const ENETUNREACH = 101    /* Network is unreachable */
const ENETRESET = 102      /* Network dropped connection because of reset */
const ECONNABORTED = 103   /* Software caused connection abort */
const ECONNRESET = 104     /* Connection reset by peer */
const ENOBUFS = 105        /* No buffer space available */
const EISCONN = 106        /* Transport endpoint is already connected */
const ENOTCONN = 107       /* Transport endpoint is not connected */
const ESHUTDOWN = 108      /* Cannot send after transport endpoint shutdown */
const ETOOMANYREFS = 109   /* Too many references: cannot splice */
const ETIMEDOUT = 110      /* Connection timed out */
const ECONNREFUSED = 111   /* Connection refused */
const EHOSTDOWN = 112      /* Host is down */
const EHOSTUNREACH = 113   /* No route to host */
const EALREADY = 114       /* Operation already in progress */
const EINPROGRESS = 115    /* Operation now in progress */
const ESTALE = 116         /* Stale NFS file handle */
const EUCLEAN = 117        /* Structure needs cleaning */
const ENOTNAM = 118        /* Not a XENIX named type file */
const ENAVAIL = 119        /* No XENIX semaphores available */
const EISNAM = 120         /* Is a named type file */
const EREMOTEIO = 121      /* Remote I/O error */
const EDQUOT = 122         /* Quota exceeded */

const ENOMEDIUM = 123    /* No medium found */
const EMEDIUMTYPE = 124  /* Wrong medium type */
const ECANCELED = 125    /* Operation Canceled */
const ENOKEY = 126       /* Required key not available */
const EKEYEXPIRED = 127  /* Key has expired */
const EKEYREVOKED = 128  /* Key has been revoked */
const EKEYREJECTED = 129 /* Key was rejected by service */

/* for robust mutexes */
const EOWNERDEAD = 130      /* Owner died */
const ENOTRECOVERABLE = 131 /* State not recoverable */

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
func Strerror(errno int32) []byte {
	b, ok := err2bytes[int(errno)]
	if ok {
		return b
	}
	return err2bytes[0]
}

// thread safe holder for the current errno
var currentErrno = NewSafe(0)

func setCurrentErrno(err error) {
	if err == nil {
		currentErrno.Set(0)
	} else {
		errStr := strings.ToLower(err.Error())
		if i, ok := errInverse[errStr]; ok {
			currentErrno.Set(i)
		} else {
			currentErrno.Set(ENODATA)
		}
	}
}

// Errno returns a pointer to the current errno.
func Errno() []int32 {
	i := currentErrno.Get().(int)
	return []int32{int32(i)}
}
