//go:build eosio
// +build eosio

package runtime

/*
#include <stdint.h>
void  eosio_assert_message( uint32_t test, const char* msg, uint32_t msg_len );
*/
import "C"

import (
	"unsafe"
)

type stringHeader struct {
	data unsafe.Pointer
	len  uintptr
}

type RevertFunction func(errMsg string)

var gRevertFn RevertFunction

func SetRevertOnPanicFn(fn RevertFunction) {
	gRevertFn = fn
}

func GetRevertOnPanicFn() RevertFunction {
	return gRevertFn
}

//Aborts processing of this action and unwinds all pending changes if the test condition is true
func Assert(test bool, msg string) {
	if !test && gRevertFn != nil {
		gRevertFn(msg)
		return
	}
	_test := uint32(0)
	if test {
		_test = 1
	}
	_msg := (*stringHeader)(unsafe.Pointer(&msg))
	C.eosio_assert_message(_test, (*C.char)(_msg.data), C.uint32_t(len(msg)))
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
		Assert(false, msg)
	case error:
		msg := "panic:" + v.Error()
		Assert(false, msg)
	default:
		Assert(false, "panic")
	}
	// printstring("panic: ")
	// printitf(message)
	// printnl()
	// eosio_assert(false, msg)
}

// Cause a runtime panic, which is (currently) always a string.
func runtimePanic(msg string) {
	Assert(false, "panic: runtime error: "+msg)
	// printstring("panic: runtime error: ")
	// println(msg)
	// abort()
}

// Try to recover a panicking goroutine.
func _recover(useParentFrame bool) interface{} {
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

// Panic when there is a divide by zero.
func divideByZeroPanic() {
	runtimePanic("divide by zero")
}

func blockingPanic() {
	runtimePanic("trying to do blocking operation in exported function")
}

// Panic when trying to add an entry to a nil map
func nilMapPanic() {
	runtimePanic("assignment to entry in nil map")
}
