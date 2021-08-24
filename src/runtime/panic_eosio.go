// +build eosio

package runtime

/*
	#include <stdint.h>
	void eosio_assert(int32_t test, const char* msg);
*/
import "C"

func eosio_assert(test bool, msg string) {
	buf := Alloc(uintptr(len(msg) + 1))
	pp := (*[1 << 30]byte)(buf)
	copy(pp[:], msg)
	pp[len(msg)] = 0
	if !test {
		C.eosio_assert(0, (*C.char)(buf))
	}
}

// trap is a compiler hint that this function cannot be executed. It is
// translated into either a trap instruction or a call to abort().
//export llvm.trap
func trap()

// Builtin function panic(msg), used as a compiler intrinsic.
func _panic(message interface{}) {
	switch v := message.(type) {
	case string:
		msg := "panic: " + v
		eosio_assert(false, msg)
	case error:
		msg := "panic:" + v.Error()
		eosio_assert(false, msg)
	default:
		eosio_assert(false, "panic")
	}
	// printstring("panic: ")
	// printitf(message)
	// printnl()
	// eosio_assert(false, msg)
}

// Cause a runtime panic, which is (currently) always a string.
func runtimePanic(msg string) {
	// printstring("panic: runtime error: ")
	// println(msg)
	// abort()
}

// Try to recover a panicking goroutine.
func _recover() interface{} {
	// Deferred functions are currently not executed during panic, so there is
	// no way this can return anything besides nil.
	return nil
}

// Panic when trying to dereference a nil pointer.
func nilPanic() {
	runtimePanic("nil pointer dereference")
}

// Panic when trying to acces an array or slice out of bounds.
func lookupPanic() {
	runtimePanic("index out of range")
}

// Panic when trying to slice a slice out of bounds.
func slicePanic() {
	runtimePanic("slice out of range")
}

// Panic when trying to create a new channel that is too big.
func chanMakePanic() {
	runtimePanic("new channel is too big")
}

// Panic when a shift value is negative.
func negativeShiftPanic() {
	runtimePanic("negative shift")
}

func blockingPanic() {
	runtimePanic("trying to do blocking operation in exported function")
}
