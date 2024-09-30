package book

import (
	"context"
	"github.com/tetratelabs/wazero/api"
	"reflect"
	"runtime"
	"unsafe"
)

func CopyStringToWasm(ctx context.Context, m api.Module, s string) uint64 {
	malloc := m.ExportedFunction("allocate")
	sLen := len(s)
	res, err := malloc.Call(ctx, uint64(sLen))
	if err != nil {
		panic(err)
	}
	resPtr := uint32(res[0])

	// Avoid copy to bytes just to write to wasm
	sHdr := (*reflect.StringHeader)(unsafe.Pointer(&s))
	var buf []byte
	bufHdr := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
	bufHdr.Data = sHdr.Data
	bufHdr.Len = sHdr.Len
	bufHdr.Cap = sHdr.Len
	runtime.KeepAlive(s)

	_ = m.Memory().Write(resPtr, buf)

	return (uint64(resPtr) << uint64(32)) | uint64(len(buf))
}

func CopyStringFromWasm(ctx context.Context, m api.Module, ptr, len uint32) string {
	desc, _ := m.Memory().Read(ptr, len)
	//free := m.ExportedFunction("deallocate")
	//_, err := free.Call(ctx, uint64(ptr), uint64(len))
	//if err != nil {
	//	panic(err)
	//}
	return string(desc)
}
