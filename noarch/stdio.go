package noarch

import (
	"fmt"
	"io"
	"io/ioutil"
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
	var file *os.File
	var err error

	// TODO: There are lots of other file modes to handle.
	switch mode {
	case "r":
		file, err = os.Open(filePath)
	case "w":
		file, err = os.Create(filePath)
	default:
		panic(fmt.Sprintf("unsupported file mode: %s", mode))
	}

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
	// Be senstive to NULL-terminated strings.
	length := 0
	for _, b := range []byte(content) {
		if b == 0 {
			break
		}

		length++
	}

	n, err := f.RealHandle.WriteString(content[:length])
	if err != nil {
		panic(err)
	}

	return n
}

func Tmpfile() *File {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		return nil
	}

	return &File{
		RealHandle: f,
	}
}

func Fgets(dest string, num int, f *File) string {
	buf := make([]byte, num)
	n, err := f.RealHandle.Read(buf)

	// FIXME: Is this the right thing to do in this case?
	if err != nil {
		return ""
	}

	// TODO: The result should go onto the dest. For that to work the dest needs
	// to be passed by reference.
	//
	//     *dest = string(buf)

	// Be careful to crop the buffer to the real number of bytes read.
	return string(buf[:n])
}

func Rewind(f *File) {
	f.RealHandle.Seek(0, 0)
}

func Feof(f *File) int {
	return 0
	// FIXME: This is a really bad way of doing this. Basically try and peek
	// ahead to test for EOF.
	buf := make([]byte, 1)
	_, err := f.RealHandle.Read(buf)

	result := 0
	if err == io.EOF {
		result = 1
	}

	// Undo cursor before returning.
	f.RealHandle.Seek(-1, 1)

	return result
}

func NewFile(f *os.File) *File {
	return &File{
		RealHandle: f,
	}
}
