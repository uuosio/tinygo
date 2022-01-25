//go:build eosio
// +build eosio

package os

type syscallFd = int

func (f unixFileHandle) ReadAt(b []byte, offset int64) (int, error) {
	return 0, nil
}

func (f unixFileHandle) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

func rename(oldname, newname string) error {
	return nil
}

func fixLongPath(path string) string {
	return path
}

func Readlink(name string) (string, error) {
	return "", ErrNotImplemented
}

func tempDir() string {
	return "/tmp"
}
