package main

/*
#cgo CFLAGS : -I../include
#cgo LDFLAGS: -L../lib -ltest
#include "test.h"
*/
import "C"

func main() {
	C.test()
}