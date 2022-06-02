//go:build !wasi && !darwin && !eosio
// +build !wasi,!darwin,!eosio

package syscall

func (e Errno) Is(target error) bool { return false }
