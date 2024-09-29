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

var Book_SetDescription = api.GoModuleFunc(func(ctx context.Context, module api.Module, stack []uint64) {
	ptr := stack[0]
	descPos := uint32(stack[1])
	descLen := uint32(stack[2])

	desc, ok := module.Memory().Read(descPos, descLen)
	if !ok {
		return
	}
	if obj, ok := arena.Load(ctx, ptr).(*Book); ok {
		obj.Description = string(desc)
	}
})

var Book_GetDescription = api.GoModuleFunc(func(ctx context.Context, module api.Module, stack []uint64) {
	ptr := stack[0]
	desc := ""
	if obj, ok := arena.Load(ctx, ptr).(*Book); ok {
		desc = obj.Description
	}
	stack[0] = CopyStringToWasm(ctx, module, desc)
})
