package main

/*
#cgo LDFLAGS: ./build/test.o
void say_hello();
*/
import "C"

func main() {
	C.say_hello()
}
