package book

import (
	"context"
	"github.com/tetratelabs/wazero/api"
	"reflect"
	"runtime"
	"unsafe"
)

func CopyStringToWasm(ctx context.Context, m api.Module, s string) uint64 {
	malloc := m.ExportedFunction("malloc")
	sLen := len(s)
	res, _ := malloc.Call(ctx, uint64(sLen))
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
