package noarch

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"
)

// This definition has been translated from the original definition for __sFILE,
// which is an alias for FILE. Not all of the attributes have been translated.
// They should be turned on as needed.
type File struct {
	// This is not part of the original struct but it is needed for internal
	// calls in Go.
	OsFile *os.File

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

func Fopen(filePath, mode []byte) *File {
	var file *os.File
	var err error

	sFilePath := NullTerminatedByteSlice(filePath)

	// TODO: Only some modes are supported by fopen()
	// https://github.com/elliotchance/c2go/issues/89
	switch NullTerminatedByteSlice(mode) {
	case "r":
		file, err = os.Open(sFilePath)
	case "r+":
		file, err = os.OpenFile(sFilePath, os.O_RDWR, 0)
	case "w":
		file, err = os.Create(sFilePath)
	case "w+":
		file, err = os.OpenFile(sFilePath, os.O_RDWR|os.O_CREATE, 0)
	default:
		panic(fmt.Sprintf("unsupported file mode: %s", mode))
	}

	if err != nil {
		return nil
	}

	return NewFile(file)
}

func Fclose(f *File) int {
	err := f.OsFile.Close()
	if err != nil {
		// Is this the correct error code?
		return 1
	}

	return 0
}

func Remove(filePath []byte) int {
	if os.Remove(string(filePath)) != nil {
		return -1
	}

	return 0
}

func Rename(from, to []byte) int {
	if os.Rename(string(from), string(to)) != nil {
		return -1
	}

	return 0
}

func Fputs(content []byte, f *File) int {
	// Be senstive to NULL-terminated strings.
	length := 0
	for _, b := range []byte(content) {
		if b == 0 {
			break
		}

		length++
	}

	n, err := f.OsFile.WriteString(string(content[:length]))
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

	return NewFile(f)
}

func Fgets(dest []byte, num int, f *File) []byte {
	buf := make([]byte, num)
	n, err := f.OsFile.Read(buf)

	// FIXME: Is this the right thing to do in this case?
	if err != nil {
		return []byte{}
	}

	// TODO: Allow arguments to be passed by reference.
	// https://github.com/elliotchance/c2go/issues/90
	// This appears in multiple locations.

	// Be careful to crop the buffer to the real number of bytes read.
	//
	// We do not trim off the NULL characters because we do not know if the file
	// we are reading is in binary mode.
	if n == num {
		// If it is the case that we have read the entire buffer with this read
		// we need to make sure we leave room for what would be the NULL
		// character at the end of the string in C.
		return buf[:n-1]
	}

	return buf[:n]
}

func Rewind(f *File) {
	f.OsFile.Seek(0, 0)
}

func Feof(f *File) int {
	// FIXME: This is a really bad way of doing this. Basically try and peek
	// ahead to test for EOF.
	buf := make([]byte, 1)
	_, err := f.OsFile.Read(buf)

	result := 0
	if err == io.EOF {
		result = 1
	}

	// Undo cursor before returning.
	f.OsFile.Seek(-1, 1)

	return result
}

func NewFile(f *os.File) *File {
	return &File{
		OsFile: f,
	}
}

func Tmpnam(buffer []byte) []byte {
	// TODO: Allow arguments to be passed by reference.
	// https://github.com/elliotchance/c2go/issues/90
	// This appears in multiple locations.

	// TODO: There must be a better way of doing this. This way allows the same
	// great distinct Go temp file generation (that also checks for existing
	// files), but unfortunately creates the file in the process; even if you
	// don't intend to use it.
	f, err := ioutil.TempFile("", "")
	if err != nil {
		return []byte{}
	}

	f.Close()
	return []byte(f.Name())
}

func Fflush(f *File) int {
	err := f.OsFile.Sync()
	if err != nil {
		return 1
	}

	return 0
}

func Fprintf(f *File, format []byte, args ...interface{}) int {
	n, err := fmt.Fprintf(f.OsFile, string(format), args...)
	if err != nil {
		return -1
	}

	return n
}

func Fscanf(f *File, format []byte, args ...interface{}) int {
	n, err := fmt.Fscanf(f.OsFile, string(format), args...)
	if err != nil {
		return -1
	}

	return n
}

func getc(f *os.File) int {
	buffer := make([]byte, 1)
	_, err := f.Read(buffer)
	if err != nil {
		return -1
	}

	return int(buffer[0])
}

func Fgetc(f *File) int {
	return getc(f.OsFile)
}

func Fputc(c int, f *File) int {
	n, err := f.OsFile.Write([]byte{byte(c)})
	if err != nil {
		return 0
	}

	return n
}

func Getchar() int {
	return getc(os.Stdin)
}

func Fseek(f *File, offset int32, origin int) int {
	n, err := f.OsFile.Seek(int64(offset), origin)
	if err != nil {
		return -1
	}

	return int(n)
}

func Ftell(f *File) int32 {
	return int32(Fseek(f, 0, 1))
}

func Fread(buffer *[]byte, size1, size2 int, f *File) int {
	// Create a new buffer so that we can ensure we read up to the correct
	// number of bytes from the file.
	newBuffer := make([]byte, size1*size2)
	n, err := f.OsFile.Read(newBuffer)

	// Despite any error we need to make sure the bytes read are copied to the
	// destination buffer.
	for i, b := range newBuffer {
		(*buffer)[i] = b
	}

	// Now we can handle the success or failure.
	if err != nil {
		return -1
	}

	return n
}

func Fwrite(buffer []byte, size1, size2 int, f *File) int {
	n, err := f.OsFile.Write(buffer[:size1*size2])
	if err != nil {
		return -1
	}

	return n
}

func Fgetpos(f *File, pos *int) int {
	absolutePos := Fseek(f, 0, 1)
	if pos != nil {
		*pos = absolutePos
	}

	return absolutePos
}

func Fsetpos(f *File, pos *int) int {
	return Fseek(f, int32(*pos), 0)
}

func Printf(format []byte, args ...interface{}) int {
	realArgs := []interface{}{}

	// Convert any C strings into Go strings.
	typeOfByteSlice := reflect.TypeOf([]byte(nil))
	for _, arg := range args {
		if reflect.TypeOf(arg) == typeOfByteSlice {
			realArgs = append(realArgs, NullTerminatedByteSlice(arg.([]byte)))
		} else {
			realArgs = append(realArgs, arg)
		}
	}

	n, _ := fmt.Printf(NullTerminatedByteSlice(format), realArgs...)

	return n
}

func Puts(s []byte) int {
	n, _ := fmt.Println(NullTerminatedByteSlice(s))

	return n
}

func Scanf(format []byte, args ...interface{}) int {
	n, _ := fmt.Scanf(NullTerminatedByteSlice(format), args...)

	return n
}

func Putchar(character int) {
	fmt.Printf("%c", character)
}
