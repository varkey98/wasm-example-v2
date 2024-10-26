package main

import (
	"context"
	_ "embed"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	"os"
	"sync"
	"wasm-example-v2/book"
)

//go:embed process.wasm
var bytes []byte

type Processor struct {
	module     wazero.CompiledModule
	cfg        wazero.ModuleConfig
	modulePool sync.Pool
	r          wazero.Runtime
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
		WithStderr(os.Stderr).
		WithName("")
	if err != nil {
		return nil, err
	}

	return &Processor{
		module: code,
		cfg:    config,
		r:      r,
	}, nil
}

func (p *Processor) GetModule() api.Module {
	return p.modulePool.Get().(api.Module)
}

func (p *Processor) ResetModule(module api.Module) {
	p.modulePool.Put(module)
}
func (p *Processor) Close(ctx context.Context) error {
	return p.module.Close(ctx)
}
