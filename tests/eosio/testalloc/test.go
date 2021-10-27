package main

import (
	"runtime"

	"github.com/uuosio/chain"
)

func main() {
	a := uint64(uintptr(runtime.Alloc(1)))
	b := uint64(uintptr(runtime.Alloc(1)))
	chain.Println("++++", a, b)
	chain.Check(a > 8192, "bad alloc start pointer")
	chain.Check(a+8 == b, "bad alloc")
}
