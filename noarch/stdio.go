package noarch

import (
	"os"
)

// This definition has been translated from the original definition for __sFILE,
// which is an alias for FILE. Not all of the attributes have been translated.
// They should be turned on as needed.
type File struct {
	// This is not part of the original struct but it is needed for internal
	// calls in Go.
	RealHandle *os.File

	// unsigned char *_p;
	// int _r;
	// int _w;
	// short _flags;
	// short _file;
	// struct __sbuf _bf;
	// int _lbfsize;
	// void *_cookie;
	// int (* _Nullable _close)(void *);
	// int (* _Nullable _read) (void *, char *, int);
	// fpos_t (* _Nullable _seek) (void *, fpos_t, int);
	// int (* _Nullable _write)(void *, const char *, int);
	// struct __sbuf _ub;
	// struct __sFILEX *_extra;
	// int _ur;
	// unsigned char _ubuf[3];
	// unsigned char _nbuf[1];
	// struct __sbuf _lb;
	// int _blksize;
	// fpos_t _offset;
}

func Fopen(filePath, mode string) *File {
	file, err := os.Open(filePath)
	if err != nil {
		return nil
	}

	return &File{
		RealHandle: file,
	}
}

func Fclose(f *File) int {
	err := f.RealHandle.Close()
	if err != nil {
		// Is this the correct error code?
		return 1
	}

	return 0
}

func Remove(filePath string) int {
	if os.Remove(filePath) != nil {
		return -1
	}

	return 0
}

func Rename(from, to string) int {
	if os.Rename(from, to) != nil {
		return -1
	}

	return 0
}

func Fputs(content string, f *File) int {
	n, _ := f.RealHandle.WriteString(content)

	return n
}
