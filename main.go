package main

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/google/uuid"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	"log"
	"os"
	"time"
	"wasm-example-v2/arena"
	"wasm-example-v2/book"
)

//go:embed body.json
var description []byte

const GoRoutines = 501

func main() {
	start := time.Now()
	MultipleGoRoutinesTest()
	fmt.Printf("Elapsed Time for WASM based invocation: %v\n", time.Since(start))

}

func MultipleGoRoutinesTest() {
	ctx := context.Background()
	ctx = arena.Initialize(ctx)
	req := book.Book{
		Name:        uuid.New().String(),
		Description: string(description),
	}

	cc, err := wazero.NewCompilationCacheWithDir("/Users/varkeychanjacob/Projects/wasm-example-v2/wazerocache")
	if err != nil {
		log.Fatal(err)
	}
	r := wazero.NewRuntimeWithConfig(ctx, wazero.NewRuntimeConfig().WithCompilationCache(cc))
	hostBuilder := r.NewHostModuleBuilder("env")
	book.AddToModule(hostBuilder)
	_, err = hostBuilder.Instantiate(ctx)
	if err != nil {
		fmt.Println(err)
	}
	_, err = wasi_snapshot_preview1.Instantiate(ctx, r)
	if err != nil {
		fmt.Println(err)
	}
	config := wazero.NewModuleConfig().
		WithStdout(os.Stdout).
		WithStderr(os.Stderr).
		WithName("")

	compiledMod1 := getCompiledModule(r, ctx)
	module1, err := r.InstantiateModule(context.Background(), compiledMod1, config)
	if err != nil {
		fmt.Println(err)
		return
	}
	fn := module1.ExportedFunction("Process")
	if fn == nil {
		fmt.Printf(" no function named Process")
		return
	}
	res, err := fn.Call(ctx, arena.Store(ctx, &req))
	if err != nil {
		log.Panicf("failed to invoke Process: %v", err)
	}

	if obj, ok := arena.Load(ctx, res[0]).(*book.Book); ok {
		fmt.Printf("Processed Value: %s\n", obj.Name)
	}

	compiledMod2 := getCompiledModule(r, ctx)

	module1.Close(ctx)
	compiledMod1.Close(ctx)
	module2, err := r.InstantiateModule(ctx, compiledMod2, config)
	if err != nil {
		fmt.Println(err)
		return
	}

	fn = module2.ExportedFunction("Process")
	if fn == nil {
		fmt.Printf(" no function named Process")
		return
	}

	res, err = fn.Call(ctx, arena.Store(ctx, &req))
	if err != nil {
		log.Panicf("failed to invoke Process: %v", err)
	}

	if obj, ok := arena.Load(ctx, res[0]).(*book.Book); ok {
		fmt.Printf("Processed Value: %s\n", obj.Name)
	}

	module2.Close(context.Background())
}

func getCompiledModule(r wazero.Runtime, ctx context.Context) wazero.CompiledModule {
	code, err := r.CompileModule(ctx, bytes)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return code
}
