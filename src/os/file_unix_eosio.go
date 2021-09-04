// +build eosio

package os

/*
#include <stdint.h>
void prints_l( const char* cstr, uint32_t len);
*/
import "C"
import (
	"errors"
	"syscall"
	"unsafe"
)

func PrintBytes(s []byte) {
	C.prints_l((*C.char)(unsafe.Pointer(&s[0])), C.uint32_t(len(s)))
}

func init() {
	// Mount the host filesystem at the root directory. This is what most
	// programs will be expecting.
	Mount("/", unixFilesystem{})
}

// Stdin, Stdout, and Stderr are open Files pointing to the standard input,
// standard output, and standard error file descriptors.
var (
	Stdin  = &File{unixFileHandle(0), "/dev/stdin"}
	Stdout = &File{unixFileHandle(1), "/dev/stdout"}
	Stderr = &File{unixFileHandle(2), "/dev/stderr"}
)

// isOS indicates whether we're running on a real operating system with
// filesystem support.
const isOS = true

// unixFilesystem is an empty handle for a Unix/Linux filesystem. All operations
// are relative to the current working directory.
type unixFilesystem struct {
}

func (fs unixFilesystem) Mkdir(path string, perm FileMode) error {
	return errors.New("unimplemented")
}

func (fs unixFilesystem) Remove(path string) error {
	return errors.New("unimplemented")
}

func (fs unixFilesystem) OpenFile(path string, flag int, perm FileMode) (FileHandle, error) {
	// Map os package flags to syscall flags.
	return unixFileHandle(0), errors.New("unimplemented")
}

// unixFileHandle is a Unix file pointer with associated methods that implement
// the FileHandle interface.
type unixFileHandle uintptr

// Read reads up to len(b) bytes from the File. It returns the number of bytes
// read and any error encountered. At end of file, Read returns 0, io.EOF.
func (f unixFileHandle) Read(b []byte) (n int, err error) {
	return 0, errors.New("unimplemented")
}

// Write writes len(b) bytes to the File. It returns the number of bytes written
// and an error, if any. Write returns a non-nil error when n != len(b).
func (f unixFileHandle) Write(b []byte) (n int, err error) {
	PrintBytes(b)
	n = len(b)
	err = nil
	return n, err
}

// Close closes the File, rendering it unusable for I/O.
func (f unixFileHandle) Close() error {
	return errors.New("unimplemented")
}

// handleSyscallError converts syscall errors into regular os package errors.
// The err parameter must be either nil or of type syscall.Errno.
func handleSyscallError(err error) error {
	if err == nil {
		return nil
	}
	switch err.(syscall.Errno) {
	case syscall.EEXIST:
		return ErrExist
	case syscall.ENOENT:
		return ErrNotExist
	default:
		return err
	}
}
