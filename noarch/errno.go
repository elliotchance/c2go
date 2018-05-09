package noarch

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

var err2str = map[int]string{
	EPERM:   "Operation not permitted",
	ENOENT:  "No such file or directory",
	ESRCH:   "No such process",
	EINTR:   "Interrupted system call",
	EIO:     "I/O error",
	ENXIO:   "No such device or address",
	E2BIG:   "Argument list too long",
	ENOEXEC: "Exec format error",
	EBADF:   "Bad file number",
	ECHILD:  "No child processes",
	EAGAIN:  "Try again",
	ENOMEM:  "Out of memory",
	EACCES:  "Permission denied",
	EFAULT:  "Bad address",
	ENOTBLK: "Block device required",
	EBUSY:   "Device or resource busy",
	EEXIST:  "File exists",
	EXDEV:   "Cross-device link",
	ENODEV:  "No such device",
	ENOTDIR: "Not a directory",
	EISDIR:  "Is a directory",
	EINVAL:  "Invalid argument",
	ENFILE:  "File table overflow",
	EMFILE:  "Too many open files",
	ENOTTY:  "Not a typewriter",
	ETXTBSY: "Text file busy",
	EFBIG:   "File too large",
	ENOSPC:  "No space left on device",
	ESPIPE:  "Illegal seek",
	EROFS:   "Read-only file system",
	EMLINK:  "Too many links",
	EPIPE:   "Broken pipe",
	EDOM:    "Math argument out of domain of func",
	ERANGE:  "Math result not representable",

	EDEADLK:      "Resource deadlock would occur",
	ENAMETOOLONG: "File name too long",
	ENOLCK:       "No record locks available",
	ENOSYS:       "Function not implemented",
	ENOTEMPTY:    "Directory not empty",
	ELOOP:        "Too many symbolic links encountered",
	//EWOULDBLOCK: "Operation would block", // same key as EAGAIN
	ENOMSG:   "No message of desired type",
	EIDRM:    "Identifier removed",
	ECHRNG:   "Channel number out of range",
	EL2NSYNC: "Level 2 not synchronized",
	EL3HLT:   "Level 3 halted",
	EL3RST:   "Level 3 reset",
	ELNRNG:   "Link number out of range",
	EUNATCH:  "Protocol driver not attached",
	ENOCSI:   "No CSI structure available",
	EL2HLT:   "Level 2 halted",
	EBADE:    "Invalid exchange",
	EBADR:    "Invalid request descriptor",
	EXFULL:   "Exchange full",
	ENOANO:   "No anode",
	EBADRQC:  "Invalid request code",
	EBADSLT:  "Invalid slot",

	EBFONT:          "Bad font file format",
	ENOSTR:          "Device not a stream",
	ENODATA:         "No data available",
	ETIME:           "Timer expired",
	ENOSR:           "Out of streams resources",
	ENONET:          "Machine is not on the network",
	ENOPKG:          "Package not installed",
	EREMOTE:         "Object is remote",
	ENOLINK:         "Link has been severed",
	EADV:            "Advertise error",
	ESRMNT:          "Srmount error",
	ECOMM:           "Communication error on send",
	EPROTO:          "Protocol error",
	EMULTIHOP:       "Multihop attempted",
	EDOTDOT:         "RFS specific error",
	EBADMSG:         "Not a data message",
	EOVERFLOW:       "Value too large for defined data type",
	ENOTUNIQ:        "Name not unique on network",
	EBADFD:          "File descriptor in bad state",
	EREMCHG:         "Remote address changed",
	ELIBACC:         "Can not access a needed shared library",
	ELIBBAD:         "Accessing a corrupted shared library",
	ELIBSCN:         ".lib section in a.out corrupted",
	ELIBMAX:         "Attempting to link in too many shared libraries",
	ELIBEXEC:        "Cannot exec a shared library directly",
	EILSEQ:          "Illegal byte sequence",
	ERESTART:        "Interrupted system call should be restarted",
	ESTRPIPE:        "Streams pipe error",
	EUSERS:          "Too many users",
	ENOTSOCK:        "Socket operation on non-socket",
	EDESTADDRREQ:    "Destination address required",
	EMSGSIZE:        "Message too long",
	EPROTOTYPE:      "Protocol wrong type for socket",
	ENOPROTOOPT:     "Protocol not available",
	EPROTONOSUPPORT: "Protocol not supported",
	ESOCKTNOSUPPORT: "Socket type not supported",
	EOPNOTSUPP:      "Operation not supported on transport endpoint",
	EPFNOSUPPORT:    "Protocol family not supported",
	EAFNOSUPPORT:    "Address family not supported by protocol",
	EADDRINUSE:      "Address already in use",
	EADDRNOTAVAIL:   "Cannot assign requested address",
	ENETDOWN:        "Network is down",
	ENETUNREACH:     "Network is unreachable",
	ENETRESET:       "Network dropped connection because of reset",
	ECONNABORTED:    "Software caused connection abort",
	ECONNRESET:      "Connection reset by peer",
	ENOBUFS:         "No buffer space available",
	EISCONN:         "Transport endpoint is already connected",
	ENOTCONN:        "Transport endpoint is not connected",
	ESHUTDOWN:       "Cannot send after transport endpoint shutdown",
	ETOOMANYREFS:    "Too many references: cannot splice",
	ETIMEDOUT:       "Connection timed out",
	ECONNREFUSED:    "Connection refused",
	EHOSTDOWN:       "Host is down",
	EHOSTUNREACH:    "No route to host",
	EALREADY:        "Operation already in progress",
	EINPROGRESS:     "Operation now in progress",
	ESTALE:          "Stale NFS file handle",
	EUCLEAN:         "Structure needs cleaning",
	ENOTNAM:         "Not a XENIX named type file",
	ENAVAIL:         "No XENIX semaphores available",
	EISNAM:          "Is a named type file",
	EREMOTEIO:       "Remote I/O error",
	EDQUOT:          "Quota exceeded",

	ENOMEDIUM:    "No medium found",
	EMEDIUMTYPE:  "Wrong medium type",
	ECANCELED:    "Operation Canceled",
	ENOKEY:       "Required key not available",
	EKEYEXPIRED:  "Key has expired",
	EKEYREVOKED:  "Key has been revoked",
	EKEYREJECTED: "Key was rejected by service",

	// for robust mutexes
	EOWNERDEAD:      "Owner died",
	ENOTRECOVERABLE: "State not recoverable",
}

var err2bytes map[int][]byte

func init() {
	err2bytes[0] = []byte{0}
	for i, str := range err2str {
		err2bytes[i] = []byte(str)
		err2bytes[i] = append(err2bytes[i], 0)
	}
}

func Strerror(errno int32) []byte {
	b, ok := err2bytes[int(errno)]
	if ok {
		return b
	}
	return err2bytes[0]
}
