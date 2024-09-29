package book

import (
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

func AddToModule(mod wazero.HostModuleBuilder) {
	mod.
		NewFunctionBuilder().
		WithGoModuleFunction(
			Book_SetName,
			[]api.ValueType{
				api.ValueTypeI64, // Ptr
				api.ValueTypeI32, // Name position
				api.ValueTypeI32, // Name length
			},
			[]api.ValueType{}).
		Export("Book_SetName").
		NewFunctionBuilder().
		WithGoModuleFunction(
			Book_GetName,
			[]api.ValueType{
				api.ValueTypeI64, // Ptr
			},
			[]api.ValueType{
				api.ValueTypeI64, // Name Ptr
			}).
		Export("Book_GetName").
		NewFunctionBuilder().
		WithGoModuleFunction(
			Book_SetDescription,
			[]api.ValueType{
				api.ValueTypeI64, // Ptr
				api.ValueTypeI32, // Description position
				api.ValueTypeI32, // Description length
			}, []api.ValueType{}).
		Export("Book_SetDescription").
		NewFunctionBuilder().
		WithGoModuleFunction(
			Book_GetDescription,
			[]api.ValueType{
				api.ValueTypeI64, // Ptr
			},
			[]api.ValueType{
				api.ValueTypeI64, // Description ptr
			}).
		Export("Book_GetDescription")
}
