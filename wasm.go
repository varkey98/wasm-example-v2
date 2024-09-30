//go:build wasi

package main

import (
	"reflect"
	"regexp"
	"unsafe"
)

//export free
func free(ptr uint32)

//export Book_SetName
func SetName(spanPtr uint64, kPtr, kLen uint32)

func setNameWrapper(spanPtr uint64, key string) {
	keyBytes := []byte(key)
	kPtr := toUint32(keyBytes)

	SetName(spanPtr, kPtr, uint32(len(keyBytes)))
}

//export Book_GetName
func GetName(spanPtr uint64) uint64

func getNameWrapper(spanPtr uint64) string {
	out := GetName(spanPtr)
	vLen := uint32(out)
	vPtr := uint32(out >> uint64(32))
	valBytes := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(vPtr),
		Len:  int(vLen),
		Cap:  int(vLen),
	}))
	ret := string(valBytes)
	free(vPtr)
	return ret
}

//export Book_SetDescription
func SetDescription(spanPtr uint64, kPtr, kLen uint32)

func setDescriptionWrapper(spanPtr uint64, key string) {
	keyBytes := []byte(key)
	kPtr := toUint32(keyBytes)

	SetDescription(spanPtr, kPtr, uint32(len(keyBytes)))
}

//export Book_GetDescription
func GetDescription(spanPtr uint64) uint64

func getDescriptionWrapper(spanPtr uint64) string {
	out := GetDescription(spanPtr)
	vLen := uint32(out)
	vPtr := uint32(out >> uint64(32))
	valBytes := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(vPtr),
		Len:  int(vLen),
		Cap:  int(vLen),
	}))
	ret := string(valBytes)
	free(vPtr)
	return ret
}

// alivePointers maps unsafe pointers to their corresponding values so that
// they aren't collected while in external use (in WebAssembly).
var alivePointers = map[uintptr]interface{}{}

// keepaliveBuf stores a reference to the buffer and returns its pointer.
//
// Callers must invoke the exported function FnDeallocateName to free memory.
func keepaliveBuf(buf []byte) uint32 {
	ptr := &buf[0]
	unsafePtr := uintptr(unsafe.Pointer(ptr))
	alivePointers[unsafePtr] = buf
	return uint32(unsafePtr)
}

func toUint32(b []byte) uint32 {
	return uint32(uintptr(unsafe.Pointer(&b[0])))
}

// allocate makes a buffer of the given size and returns its uintptr. Once
// finished, the caller must free the memory with FnDeallocateName.
//
//export Allocate
func allocate(size uint32) uint32 {
	return keepaliveBuf(make([]byte, size))
}

// deallocate frees a uintptr returned by keepaliveBuf or allocate, allowing it
// to be garbage collected.
//
//export Deallocate
func deallocate(ptr uint32) {
	delete(alivePointers, uintptr(ptr))
}

//export Process
func processAttributesV2(ptr uint64) uint64 {
	val := getNameWrapper(ptr)
	setNameWrapper(ptr, val+" processed")

	return ptr
}

//export ProcessRegex
func processAttributesV3(ptr uint64) uint64 {
	regex := regexp.MustCompile(`.*traceable.*`)
	val := getDescriptionWrapper(ptr)
	if regex.MatchString(val) {
		//fmt.Println("Matched")
		val = val + ": processed"
		setDescriptionWrapper(ptr, val)
	}

	return ptr
}

func main() {
	//fmt.Println("Hello World")
}
