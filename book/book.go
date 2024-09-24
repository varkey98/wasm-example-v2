package book

import (
	"context"
	"github.com/tetratelabs/wazero/api"
	"wasm-example-v2/arena"
)

type Book struct {
	Name        string
	Description string
}

var Book_SetName = api.GoModuleFunc(func(ctx context.Context, module api.Module, stack []uint64) {
	ptr := stack[0]
	namePos := uint32(stack[1])
	nameLen := uint32(stack[2])

	name, ok := module.Memory().Read(namePos, nameLen)
	if !ok {
		return
	}
	if obj, ok := arena.Load(ctx, ptr).(*Book); ok {
		obj.Name = string(name)
	}
})

var Book_GetName = api.GoModuleFunc(func(ctx context.Context, module api.Module, stack []uint64) {
	ptr := stack[0]
	name := ""
	if obj, ok := arena.Load(ctx, ptr).(*Book); ok {
		name = obj.Name
	}
	stack[0] = CopyStringToWasm(ctx, module, name)
})
