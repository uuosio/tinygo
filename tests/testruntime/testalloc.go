package main

import (
	"unsafe"

	"github.com/uuosio/chain"
	"github.com/uuosio/chain/logger"
)

func main() {
	// test align
	{
		a := make([]byte, 1)
		b := make([]byte, 1)
		_a := uint64(uintptr(unsafe.Pointer(&a[0])))
		_b := uint64(uintptr(unsafe.Pointer(&b[0])))
		logger.Println(_a, _b)
		chain.Check(_a+8 == _b, "bad value")
	}
}
