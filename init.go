package main

import (
	"context"
	_ "embed"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	"os"
	"wasm-example-v2/book"
)

//go:embed process.wasm
var bytes []byte

type Processor struct {
	module api.Module
}

func Initialise() (*Processor, error) {
	ctx := context.Background()

	r := wazero.NewRuntime(ctx)

	hostBuilder := r.NewHostModuleBuilder("env")
	book.AddToModule(hostBuilder)
	_, err := hostBuilder.Instantiate(ctx)
	if err != nil {
		return nil, err
	}

	_, err = wasi_snapshot_preview1.Instantiate(ctx, r)
	if err != nil {
		return nil, err
	}

	code, err := r.CompileModule(ctx, bytes)
	if err != nil {
		return nil, err
	}
	config := wazero.NewModuleConfig().
		WithStdout(os.Stdout).
		WithStderr(os.Stderr)
	mod, err := r.InstantiateModule(ctx, code, config)
	if err != nil {
		return nil, err
	}

	return &Processor{module: mod}, nil
}
